package models

import (
	"fmt"

	"github.com/CJPotter10/sbs-drafts-api/utils"
)

type DraftToken struct {
	Roster   Roster `json:"roster"`
	CardId   string `json:"_cardId"`
	ImageUrl string `json:"_imageUrl"`
	Level    string `json:"_level"`
	OwnerId  string `json:"_ownerId"`
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

func (dt *DraftToken) ConvertToMetadata() *Metadata {
	return &Metadata{
		Description: "Our 10,000 Spoiled Banana Society Genesis Cards minted on the Ethereum blockchain doubles as your membership and gives you access to the Spoiled Banana Society benefits including playing each year in our SBS Genesis League with no further purchase necessary.",
		Name:        fmt.Sprintf("SBS Draft Token %s", dt.CardId),
		Image: 		 dt.ImageUrl,
		Attributes: CreateTokenAttributes(dt),	
	}
}

func CreateTokenAttributes(dt *DraftToken) []AttributeType {
	res := make([]AttributeType, 0) 
	for i := 0; i < len(dt.Roster.QB); i++ {
		obj := AttributeType{
			Type: fmt.Sprintf("QB%x", i),
			Value: dt.Roster.QB[i],
		}
		res = append(res, obj)
	} 
	for i := 0; i < len(dt.Roster.RB); i++ {
		obj := AttributeType{
			Type: fmt.Sprintf("RB%x", i),
			Value: dt.Roster.RB[i],
		}
		res = append(res, obj)
	} 
	for i := 0; i < len(dt.Roster.TE); i++ {
		obj := AttributeType{
			Type: fmt.Sprintf("TE%x", i),
			Value: dt.Roster.TE[i],
		}
		res = append(res, obj)
	} 
	for i := 0; i < len(dt.Roster.WR); i++ {
		obj := AttributeType{
			Type: fmt.Sprintf("WR%x", i),
			Value: dt.Roster.QB[i],
		}
		res = append(res, obj)
	} 
	for i := 0; i < len(dt.Roster.DST); i++ {
		obj := AttributeType{
			Type: fmt.Sprintf("DST%x", i),
			Value: dt.Roster.DST[i],
		}
		res = append(res, obj)
	} 

	levelTrait := AttributeType{
		Type: "LEVEL",
		Value: dt.Level,
	}
	
	res = append(res, levelTrait)

	return res
}

func (token *DraftToken) GetDraftTokenFromDraftById(tokenId, draftId string) error {
	err := utils.Db.ReadDocument(fmt.Sprintf("drafts/%s/cards", draftId), tokenId, token)
	if err != nil {
		return err
	}
	return nil
}
