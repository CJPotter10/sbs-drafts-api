package models

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/CJPotter10/sbs-drafts-api/utils"
)

type DraftInfo struct {
	DraftId           string       `json:"draftId"`
	DisplayName       string       `json:"displayName"`
	DraftStartTime    time.Time    `json:"draftStartTime"`
	CurrentDrafter    string       `json:"currentDrafter"`
	CurrentPickNumber int          `json:"pickNumber"`
	CurrentRound      int          `json:"roundNum"`
	PickInRound       int          `json:"pickInRound"`
	DraftOrder        []LeagueUser `json:"draftOrder"`
}

func CreateDraftInfoForDraft(draftId, draftType string, currentUsers []LeagueUser, leagueInfo *League) (*DraftInfo, error) {

	draftOrder := make([]LeagueUser, len(currentUsers))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(currentUsers))

	for i, v := range perm {
		draftOrder[i] = currentUsers[v]
	}

	var startTime time.Time

	if strings.ToLower(draftType) == "live" {
		startTime = time.Now().Add(1 * time.Minute)
	} else {
		res, err := findTheNextSaturday()
		if err != nil {
			return nil, err
		}
		startTime = res
	}

	res := &DraftInfo{
		DraftId:           draftId,
		DisplayName:       leagueInfo.DisplayName,
		DraftStartTime:    startTime,
		CurrentDrafter:    draftOrder[0].OwnerId,
		CurrentPickNumber: 1,
		CurrentRound:      1,
		PickInRound:       1,
		DraftOrder:        draftOrder,
	}

	return res, nil
}

func findTheNextSaturday() (time.Time, error) {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := 6
	hour := 18
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		fmt.Println("Error finding the LA timezone or location")
		return time.Time{}, err
	}

	startTime := time.Date(year, month, day, hour, 0, 0, 0, loc)
	return startTime, nil

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

func (info *DraftInfo) Update(draftId string) error {
	err := utils.Db.CreateOrUpdateDocument(fmt.Sprintf("drafts/%s/state", draftId), "info", info)
	if err != nil {
		return err
	}

	return nil
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

func CreateDraftSummaryForDraft(draftId string, draftOrder []LeagueUser) *DraftSummary {
	sum := &DraftSummary{
		Summary: make([]PlayerStateInfo, 0),
	}

	pickNum := 1

	for i := 1; i <= 15; i++ {
		round := i
		for j := 1; j <= 10; j++ {
			pickInRound := j
			var drafter string
			if round%2 == 0 {
				drafter = draftOrder[len(draftOrder)-pickInRound].OwnerId
			} else {
				drafter = draftOrder[pickInRound-1].OwnerId
			}

			obj := PlayerStateInfo{
				PickNum:      pickNum,
				OwnerAddress: drafter,
				Round:        round,
			}
			sum.Summary = append(sum.Summary, obj)
			pickNum++
		}
	}
	return sum
}

func (s *DraftSummary) Update(draftId string) error {
	err := utils.Db.CreateOrUpdateDocument(fmt.Sprintf("drafts/%s/state", draftId), "summary", s)
	if err != nil {
		return err
	}

	return nil
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

func (connList *ConnectionList) Update(draftId string) error {
	err := utils.Db.CreateOrUpdateDocument(fmt.Sprintf("drafts/%s/state", draftId), "connectionList", connList)
	if err != nil {
		return err
	}

	return nil
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
	Rosters map[string]*Roster `json:"rosters"`
}

func CreateEmptyRosterState(info DraftInfo) *RosterState {
	data := make(map[string]*Roster)

	for i := 0; i < len(info.DraftOrder); i++ {
		data[info.DraftOrder[i].OwnerId] = NewEmptyRoster()
	}

	return &RosterState{
		Rosters: data,
	}
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

func GetDefaultPlayerState() (map[string]PlayerStateInfo, error) {
	data := make(map[string]PlayerStateInfo)

	err := utils.Db.ReadDocument("playerStats2023", "defaultPlayerDraftState", &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (rs *RosterState) Update(draftId string) error {
	err := utils.Db.CreateOrUpdateDocument(fmt.Sprintf("drafts/%s/state", draftId), "rosters", rs)
	if err != nil {
		return err
	}

	return nil
}

type Players struct {
	Players map[string]PlayerStateInfo `json:"players"`
}

func CreateLeagueDraftStateUponFilling(draftId string, draftType string) error {
	var leagueInfo League
	err := utils.Db.ReadDocument("drafts", draftId, &leagueInfo)
	if err != nil {
		fmt.Println("Error in reading the league document")
		return err
	}

	var counts DraftLeagueTracker
	err = utils.Db.ReadDocument("drafts", "draftTracker", &counts)
	if err != nil {
		fmt.Println("Error in reading the draft tracker document into objects")
		return err
	}

	if s := strings.ToLower(draftType); s == "live" {
		counts.CurrentLiveDraftCount++
		counts.FilledLeaguesCount++
	} else {
		counts.CurrentScheduledDraftCount++
		counts.FilledLeaguesCount++
	}

	leagueInfo.DisplayName = fmt.Sprintf("Draft League %d", counts.FilledLeaguesCount)

	err = utils.Db.CreateOrUpdateDocument("drafts", draftId, &leagueInfo)
	if err != nil {
		return err
	}

	for i := 0; i < len(leagueInfo.CurrentUsers); i++ {
		token := DraftToken{
			Roster: NewEmptyRoster(),
		}
		err := utils.Db.ReadDocument("draftTokens", leagueInfo.CurrentUsers[i].TokenId, &token)
		if err != nil {
			return err
		}

		token.LeagueDisplayName = leagueInfo.DisplayName
		err = utils.Db.CreateOrUpdateDocument("draftTokens", leagueInfo.CurrentUsers[i].TokenId, &token)
		if err != nil {
			return err
		}
		fmt.Println("Updated display name on card ", leagueInfo.CurrentUsers[i].TokenId)
	}

	err = utils.Db.CreateOrUpdateDocument("drafts", "draftTracker", counts)
	if err != nil {
		return err
	}

	if len(leagueInfo.CurrentUsers) != 10 {
		return fmt.Errorf("there is not 10 users in this league so we can not make a draft state for an unfilled league")
	}

	info, err := CreateDraftInfoForDraft(draftId, leagueInfo.DraftType, leagueInfo.CurrentUsers, &leagueInfo)
	if err != nil {
		return err
	}
	if err := info.Update(draftId); err != nil {
		return err
	}

	data, err := GetDefaultPlayerState()
	if err != nil {
		return err
	}
	fmt.Println("Data returned from get default player state")

	err = utils.Db.CreateOrUpdateDocument(fmt.Sprintf("drafts/%s/state", draftId), "playerState", &data)
	if err != nil {
		return err
	}

	summary := CreateDraftSummaryForDraft(draftId, info.DraftOrder)
	if err := summary.Update(draftId); err != nil {
		return err
	}
	connList := CreateNewConnectionList(*info)
	if err := connList.Update(draftId); err != nil {
		return err
	}
	rosterMap := CreateEmptyRosterState(*info)
	if err := rosterMap.Update(draftId); err != nil {
		return err
	}

	return nil

}
