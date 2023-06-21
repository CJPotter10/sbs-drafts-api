package models

import (
	"context"
	"fmt"

	"github.com/CJPotter10/sbs-drafts-api/utils"
)

type DraftToken struct {
	Roster            *Roster `json:"roster"`
	DraftType         string  `json:"_draftType"`
	CardId            string  `json:"_cardId"`
	ImageUrl          string  `json:"_imageUrl"`
	Level             string  `json:"_level"`
	OwnerId           string  `json:"_ownerId"`
	LeagueId          string  `json:"_leagueId"`
	LeagueDisplayName string  `json:"_leagueDisplayName"`
	Rank              string  `json:"_rank"`
	WeekScore         string  `json:"_weekScore"`
	SeasonScore       string  `json:"_seasonScore"`
}

type UsersTokens struct {
	Available []DraftToken `json:"available"`
	Active    []DraftToken `json:"active"`
}

func ReturnAllDraftTokensForOwner(ownerId string) (*UsersTokens, error) {
	res := &UsersTokens{
		Available: make([]DraftToken, 0),
		Active:    make([]DraftToken, 0),
	}

	data, err := utils.Db.Client.Collection(fmt.Sprintf("owners/%s/validDraftTokens", ownerId)).Documents(context.Background()).GetAll()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(data); i++ {
		var token DraftToken
		data[i].DataTo(&token)
		res.Available = append(res.Available, token)
	}

	data, err = utils.Db.Client.Collection(fmt.Sprintf("owners/%s/usedDraftTokens", ownerId)).Documents(context.Background()).GetAll()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(data); i++ {
		var token DraftToken
		data[i].DataTo(&token)
		res.Active = append(res.Active, token)
	}

	return res, nil
}

type Metadata struct {
	Description string          `json:"description"`
	Name        string          `json:"name"`
	Image       string          `json:"image"`
	Attributes  []AttributeType `json:"attributes"`
}

type AttributeType struct {
	Type  string `json:"trait_type"`
	Value string `json:"value"`
}

func MintDraftTokenInDb(tokenId, ownerId string) (*DraftToken, error) {
	// tokenNum, err := strconv.Atoi(tokenId)
	// if err != nil {
	// 	return nil, err
	// }

	// res, err := utils.Contract.GetOwnerOfToken(tokenNum)
	// if strings.ToLower(ownerId) != res {
	// 	return nil, fmt.Errorf("the passed in ownerId does not match the ownerId returned from the smart contract for token %s: expected owner(%s) / actual owner(%s)", tokenId, ownerId, res)
	// }

	// can hardcode the image to the draft token image we will use before the draft has been complete
	draftToken := &DraftToken{
		Roster:            NewEmptyRoster(),
		DraftType:         "",
		CardId:            tokenId,
		ImageUrl:          "",
		Level:             "Pro",
		OwnerId:           ownerId,
		LeagueId:          "",
		LeagueDisplayName: "",
		Rank:              "N/A",
		WeekScore:         "0",
		SeasonScore:       "0",
	}

	err := utils.Db.CreateOrUpdateDocument("draftTokens", tokenId, draftToken)
	if err != nil {
		return nil, err
	}
	err = utils.Db.CreateOrUpdateDocument(fmt.Sprintf("owners/%s/validDraftTokens", ownerId), tokenId, draftToken)
	if err != nil {
		return nil, err
	}

	metadata := draftToken.ConvertToMetadata()
	err = utils.Db.CreateOrUpdateDocument("draftTokenMetadata", tokenId, metadata)
	if err != nil {
		return nil, err
	}

	return draftToken, nil
}

func (dt *DraftToken) ConvertToMetadata() *Metadata {
	return &Metadata{
		Description: "Our 10,000 Spoiled Banana Society Genesis Cards minted on the Ethereum blockchain doubles as your membership and gives you access to the Spoiled Banana Society benefits including playing each year in our SBS Genesis League with no further purchase necessary.",
		Name:        fmt.Sprintf("SBS Draft Token %s", dt.CardId),
		Image:       dt.ImageUrl,
		Attributes:  CreateTokenAttributes(dt),
	}
}

func CreateTokenAttributes(dt *DraftToken) []AttributeType {
	res := make([]AttributeType, 0)
	for i := 0; i < len(dt.Roster.QB); i++ {
		obj := AttributeType{
			Type:  fmt.Sprintf("QB%x", i),
			Value: dt.Roster.QB[i],
		}
		res = append(res, obj)
	}
	for i := 0; i < len(dt.Roster.RB); i++ {
		obj := AttributeType{
			Type:  fmt.Sprintf("RB%x", i),
			Value: dt.Roster.RB[i],
		}
		res = append(res, obj)
	}
	for i := 0; i < len(dt.Roster.TE); i++ {
		obj := AttributeType{
			Type:  fmt.Sprintf("TE%x", i),
			Value: dt.Roster.TE[i],
		}
		res = append(res, obj)
	}
	for i := 0; i < len(dt.Roster.WR); i++ {
		obj := AttributeType{
			Type:  fmt.Sprintf("WR%x", i),
			Value: dt.Roster.QB[i],
		}
		res = append(res, obj)
	}
	for i := 0; i < len(dt.Roster.DST); i++ {
		obj := AttributeType{
			Type:  fmt.Sprintf("DST%x", i),
			Value: dt.Roster.DST[i],
		}
		res = append(res, obj)
	}

	levelTrait := AttributeType{
		Type:  "LEVEL",
		Value: dt.Level,
	}
	res = append(res, levelTrait)

	weekScoreTrait := AttributeType{
		Type:  "WEEK SCORE",
		Value: dt.WeekScore,
	}
	res = append(res, weekScoreTrait)

	seasonScoreTrait := AttributeType{
		Type:  "Season Score",
		Value: dt.SeasonScore,
	}
	res = append(res, seasonScoreTrait)

	rankTrait := AttributeType{
		Type:  "RANK",
		Value: dt.Rank,
	}
	res = append(res, rankTrait)

	leagueTrait := AttributeType{
		Type:  "LEAGUE NAME",
		Value: dt.LeagueDisplayName,
	}
	res = append(res, leagueTrait)

	return res
}

func (token *DraftToken) GetDraftTokenFromDraftById(tokenId, draftId string) error {
	err := utils.Db.ReadDocument(fmt.Sprintf("drafts/%s/cards", draftId), tokenId, token)
	if err != nil {
		return err
	}
	return nil
}

func (token *DraftToken) updateInUseDraftTokenInDatabase() error {
	err := utils.Db.CreateOrUpdateDocument("draftTokens", token.CardId, token)
	if err != nil {
		return err
	}

	err = utils.Db.CreateOrUpdateDocument(fmt.Sprintf("owners/%s/usedDraftTokens", token.OwnerId), token.CardId, token)
	if err != nil {
		return err
	}
	err = utils.Db.CreateOrUpdateDocument(fmt.Sprintf("drafts/%s/cards", token.LeagueId), token.CardId, token)
	if err != nil {
		return err
	}

	metadata := token.ConvertToMetadata()
	err = utils.Db.CreateOrUpdateDocument("draftTokenMetadata", token.CardId, metadata)
	if err != nil {
		return err
	}

	return nil
}

func (token *DraftToken) RemoveTokenFromLeague() error {
	oldLeagueId := token.LeagueId
	token.LeagueId = ""
	token.DraftType = ""
	err := utils.Db.CreateOrUpdateDocument("draftTokens", token.CardId, token)
	if err != nil {
		return err
	}

	err = utils.Db.CreateOrUpdateDocument(fmt.Sprintf("owners/%s/validDraftTokens", token.OwnerId), token.CardId, token)
	if err != nil {
		return err
	}

	err = utils.Db.DeleteDocument(fmt.Sprintf("owners/%s/usedDraftTokens", token.OwnerId), token.CardId)
	if err != nil {
		fmt.Println("error when deleting document in owners")
		return err
	}

	fmt.Println(fmt.Sprintf("drafts/%s/cards/%s", oldLeagueId, token.CardId))
	err = utils.Db.DeleteDocument(fmt.Sprintf("drafts/%s/cards", oldLeagueId), token.CardId)
	if err != nil {
		fmt.Println("error when deleting token from draft league: ", err)
		return err
	}

	metadata := token.ConvertToMetadata()
	err = utils.Db.CreateOrUpdateDocument("draftTokenMetadata", token.CardId, metadata)
	if err != nil {
		return err
	}

	return nil

}
