package models

import (
	"fmt"

	"github.com/CJPotter10/sbs-drafts-api/utils"
)

type DraftInfo struct {
	DraftId           string   `json:"draftId"`
	CurrentDrafter    string   `json:"currentDrafter"`
	CurrentPickNumber int      `json:"pickNumber"`
	CurrentRound      int      `json:"roundNum"`
	PickInRound       int      `json:"pickInRound"`
	DraftOrder        []string `json:"draftOrder"`
}

func ReturnDraftInfoForDraft(draftId string) (*DraftInfo, error) {
	var info DraftInfo
	collectionString := fmt.Sprintf("drafts/%s/state", draftId)
	err := utils.Db.ReadDocument(collectionString, "info", &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

type DraftSummary struct {
	Summary []PlayerStateInfo `json:"summary"`
}

func ReturnDraftSummaryForDraft(draftId string) (*DraftSummary, error) {
	var sum DraftSummary
	collectionString := fmt.Sprintf("drafts/%s/state", draftId)
	err := utils.Db.ReadDocument(collectionString, "summary", &sum)
	if err != nil {
		return nil, err
	}

	return &sum, nil
}

type ConnectionList struct {
	List map[string]bool `json:"list"`
}

func ReturnConnectionListForDraft(draftId string) (*ConnectionList, error) {
	var cl ConnectionList
	collectionString := fmt.Sprintf("drafts/%s/state", draftId)
	err := utils.Db.ReadDocument(collectionString, "connectionList", &cl)
	if err != nil {
		return nil, err
	}

	return &cl, nil
}

type Roster struct {
	DST []string `json:"DST"`
	QB  []string `json:"QB"`
	RB  []string `json:"RB"`
	TE  []string `json:"TE"`
	WR  []string `json:"WR"`
}

func NewEmptyRoster() *Roster {
	return &Roster{
		DST: make([]string, 0),
		QB:  make([]string, 0),
		RB:  make([]string, 0),
		TE:  make([]string, 0),
		WR:  make([]string, 0),
	}
}

type RosterState struct {
	Rosters map[string]Roster `json:"rosters"`
}

func ReturnRostersForDraft(draftId string) (*RosterState, error) {
	var data RosterState
	collectionString := fmt.Sprintf("drafts/%s/state", draftId)
	err := utils.Db.ReadDocument(collectionString, "rosters", &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func CreateLeagueDraftStateUponFilling(draftId string) error {
	var leagueInfo League
	err := utils.Db.ReadDocument("drafts", draftId, &leagueInfo)
	if err != nil {
		fmt.Println("Error in reading the league document")
		return err
	}

	if len(leagueInfo.CurrentUsers) != 10 {
		return fmt.Errorf("there is not 10 users in this league so we can not make a draft state for an unfilled league")
	}

	return nil

	// 	IS NOT COMPLETE STILL WORK ON THIS AND COMPLETE IT

	// PLEASE DON'T FORGET THIS

}
