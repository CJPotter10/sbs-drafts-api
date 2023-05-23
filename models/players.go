package models

import (
	"fmt"
	"strings"

	"github.com/CJPotter10/sbs-drafts-api/utils"
)

type PlayerStateInfo struct {
	// unique player Id will probably just be the team and position such as BUFQB
	PlayerId string `json:"playerId"`
	// display name for front end
	DisplayName string `json:"displayName"`
	// team of the player
	Team string `json:"team"`
	// position of player
	Position string `json:"position"`
	// address of the user who drafted this player
	OwnerAddress string `json:"ownerAddress"`
	// number pick that this player was selected.... will default to nil in the database
	PickNum int `json:"pickNum"`
	// the round which this player was drafted in
	Round int `json:"round"`
}

type StateMap struct {
	Players map[string]PlayerStateInfo
}

type PlayerRanking struct {
	PlayerId string
	Ranking  int
	Score    int
}

type UserRankings struct {
	Rankings []PlayerRanking
}

type DraftPlayerRanking struct {
	// unique player Id will probably just be the team and position such as BUFQB
	PlayerId string `json:"playerId"`
	// holds the state object for player
	PlayerStateInfo PlayerStateInfo `json:"playerStateInfo"`
	Stats           StatsObject     `json:"stats"`
	Ranking         PlayerRanking   `json:"ranking"`
}

func CreateRankingObject(ranking PlayerRanking, stats StatsObject, info PlayerStateInfo) DraftPlayerRanking {
	return DraftPlayerRanking{
		PlayerStateInfo: info,
		Stats:           stats,
		Ranking:         ranking,
	}
}

func GetUserRankings(ownerId string) (*UserRankings, error) {
	var r UserRankings
	err := utils.Db.ReadDocument(fmt.Sprintf("owners/%s/drafts", ownerId), "rankings", &r)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "notfound") {
			err := utils.Db.ReadDocument("playerStats2023", "rankings", &r)
			if err != nil {
				return nil, err
			}

			err = utils.Db.CreateOrUpdateDocument(fmt.Sprintf("owners/%s/drafts", ownerId), "rankings", r)
			if err != nil {
				return nil, err
			}
		}
	}
	return &r, nil
}

type StatsObject struct {
	PlayerId     string `json:"playerId"`
	AverageScore int    `json:"averageScore"`
	HighestScore int    `json:"highestScore"`
	Top5Finishes int    `json:"top5Finishes"`
}

type StatsMap struct {
	PlayerStats map[string]StatsObject `json:"players"`
}

func ReturnPlayerStateWithRankings(ownerId string, draftId string) (map[string]DraftPlayerRanking, error) {
	userRankings, err := GetUserRankings(ownerId)
	if err != nil {
		return nil, err
	}

	var state StateMap
	err = utils.Db.ReadDocument(fmt.Sprintf("drafts/%s/state", draftId), "players", &state)
	if err != nil {
		return nil, err
	}

	var stats StatsMap
	err = utils.Db.ReadDocument("playerStats2023", "playerMap", &stats)
	if err != nil {
		return nil, err
	}

	res := make(map[string]DraftPlayerRanking)

	for _, rank := range userRankings.Rankings {
		res[rank.PlayerId] = CreateRankingObject(rank, stats.PlayerStats[rank.PlayerId], state.Players[rank.PlayerId])
	}

	return res, nil
}
