package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/CJPotter10/sbs-drafts-api/utils"
)

type League struct {
	LeagueId     string       `json:"leagueId"`
	DisplayName  string       `json:"displayName"`
	CurrentUsers []LeagueUser `json:"currentUsers"`
	NumPlayers   int          `json:"numPlayers"`
	MaxPlayers   int          `json:"maxPlayers"`
	StartDate    time.Time    `json:"startDate"`
	EndDate      time.Time    `json:"endDate"`
	DraftType    string       `json:"draftType"`
	Level        string       `json:"level"`
	IsLocked     bool         `json:"isFilled"`
}

type LeagueUser struct {
	OwnerId string `json:"ownerId"`
	TokenId string `json:"tokenId"`
}

type DraftLeagueTracker struct {
	CurrentLiveDraftCount      int `json:"currentLiveDraftCount"`
	CurrentScheduledDraftCount int `json:"currentScheduledDraftCount"`
}

func CreateLeague(ownerId string, draftNum int, draftType string) (*League, error) {
	loc, err := time.LoadLocation("America/Los_Angoles")
	if err != nil {
		fmt.Println("Error finding the chicago timezone or location")
		return nil, err
	}
	res := &League{
		LeagueId:     fmt.Sprintf("%s-draft-%x", draftType, draftNum),
		DisplayName:  fmt.Sprintf("SBS %s Draft League #%x", draftType, draftNum),
		CurrentUsers: make([]LeagueUser, 0),
		NumPlayers:   1,
		MaxPlayers:   10,
		StartDate:    time.Date(2023, time.September, 3, 0, 0, 0, 0, loc),
		EndDate:      time.Date(2023, time.December, 25, 0, 0, 0, 0, loc),
		DraftType:    draftType,
		Level:        "Pro",
		IsLocked:     false,
	}

	return res, nil
}

func JoinLeagues(ownerId string, numLeaguesToJoin int, draftType string) ([]DraftToken, error) {
	data, err := utils.Db.Client.Collection(fmt.Sprintf("owners/%s/validDraftTokens", ownerId)).Documents(context.Background()).GetAll()
	if err != nil {
		return nil, err
	}

	if len(data) < numLeaguesToJoin {
		err := fmt.Errorf("there does not seem to be enough valid draft tokens needed to enter into this number of leagues: You have %x / %x valid tokens", len(data), numLeaguesToJoin)
		return nil, err
	}

	// read document from db that tracks the amount of filled draft leagues there are for each type
	var counts DraftLeagueTracker
	err = utils.Db.ReadDocument("drafts", "draftTracker", &counts)
	if err != nil {
		fmt.Println("Error in reading the draft tracker document into objects")
		return nil, err
	}

	var currentDraft int
	if s := strings.ToLower(draftType); s == "live" {
		currentDraft = counts.CurrentLiveDraftCount
	} else {
		currentDraft = counts.CurrentScheduledDraftCount
	}

	res := make([]DraftToken, 0)

	for i := 0; i < numLeaguesToJoin; i++ {
		var t DraftToken
		err := data[i].DataTo(&t)
		if err != nil {
			return nil, err
		}
		currentDraft, err = AddCardToLeague(&t, currentDraft, draftType)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}

	return res, nil
}

func AddCardToLeague(token *DraftToken, expectedDraftNum int, draftType string) (int, error) {
	currentDraftNum := expectedDraftNum
	var draftId string
	var l League

	// find the right league to add the card to ensuring that this owner does not already have a token in that league
	for {
		draftId = fmt.Sprintf("%s-draft-%x", draftType, currentDraftNum)
		err := utils.Db.ReadDocument("drafts", draftId, &l)
		if err != nil {
			s := err.Error()
			if res := strings.Contains(s, "code = NotFound"); res {
				league, err := CreateLeague(token.OwnerId, currentDraftNum, draftType)
				if err != nil {
					return -1, err
				}
				l = *league
				break
			}
			return -1, err
		}

		isValid := true
		for j := 0; j < len(l.CurrentUsers); j++ {
			if l.CurrentUsers[j].OwnerId == token.OwnerId {
				isValid = false
			}
		}

		if isValid {
			break
		}
		currentDraftNum++
	}

	token.LeagueId = draftId
	token.DraftType = draftType

	l.CurrentUsers = append(l.CurrentUsers, LeagueUser{OwnerId: token.OwnerId, TokenId: token.CardId})
	l.NumPlayers++
	err := utils.Db.CreateOrUpdateDocument("drafts", draftId, l)
	if err != nil {
		return -1, err
	}

	// add card to league
	err = token.updateInUseDraftTokenInDatabase()
	if err != nil {
		return -1, err
	}

	_, err = utils.Db.Client.Collection(fmt.Sprintf("owners/%s/validDraftTokens", token.OwnerId)).Doc(token.CardId).Delete(context.Background())
	if err != nil {
		return -1, err
	}

	return currentDraftNum, nil
}

func RemoveUserFromDraft(tokenId string, ownerId string, draftId string) (bool, error) {
	var l League
	err := utils.Db.ReadDocument("drafts", draftId, &l)
	if err != nil {
		return false, err
	}

	if l.IsLocked || l.NumPlayers == 10 {
		return false, fmt.Errorf("this draft league is already locked so you can not leave")
	}

	isInLeague := false
	newCurrentUsers := make([]LeagueUser, 0)
	for i := 0; i < len(l.CurrentUsers); i++ {
		if l.CurrentUsers[i].OwnerId == ownerId && l.CurrentUsers[i].TokenId == tokenId {
			isInLeague = true
		} else {
			newCurrentUsers = append(newCurrentUsers, l.CurrentUsers[i])
		}
	}
	if !isInLeague {
		return false, fmt.Errorf("this user was not found to be in the current User array of the draft league")
	}

	l.CurrentUsers = newCurrentUsers
	l.NumPlayers--

	err = utils.Db.CreateOrUpdateDocument("drafts", draftId, l)
	if err != nil {
		return false, err
	}

	return true, nil
}
