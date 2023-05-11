package models

import (
	"fmt"

	"github.com/CJPotter10/sbs-drafts-api/utils"
)

type DraftInfo struct {
	DraftId           string       `json:"draftId"`
	CurrentDrafter    string       `json:"currentDrafter"`
	CurrentPickNumber int          `json:"pickNumber"`
	CurrentRound      int          `json:"roundNum"`
	PickInRound       int          `json:"pickInRound"`
	DraftOrder        []LeagueUser `json:"draftOrder"`
}

func CreateDraftInfoForDraft(draftId string, currentUsers []LeagueUser) (*DraftInfo, error) {
	// draftOrder := make([]string, 10)

	// for i := 0; i < len(currentUsers); i++ {
	// 	draftOrder[i] = currentUsers[i].OwnerId
	// }

	res := &DraftInfo{
		DraftId:           draftId,
		CurrentDrafter:    currentUsers[0].OwnerId,
		CurrentPickNumber: 1,
		CurrentRound:      1,
		PickInRound:       1,
		DraftOrder:        currentUsers,
	}

	return res, nil
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

func CreateDraftSummaryForDraft(draftId string) *DraftSummary {
	return &DraftSummary{
		Summary: make([]PlayerStateInfo, 0),
	}
}

type ConnectionList struct {
	List map[string]bool `json:"list"`
}

func CreateNewConnectionList(info DraftInfo) *ConnectionList {
	res := make(map[string]bool)
	for i := 0; i < len(info.DraftOrder); i++ {
		res[info.DraftOrder[i].OwnerId] = false
	}

	return &ConnectionList{
		List: res,
	}
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

func NewEmptyRoster() Roster {
	return Roster{
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

func 

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

	info, err := CreateDraftInfoForDraft(draftId, leagueInfo.CurrentUsers)
	if err != nil {
		return err
	}

	summary := CreateDraftSummaryForDraft(draftId)
	connList := CreateNewConnectionList(*info)



	return nil

}
