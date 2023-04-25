package models

import (
	"fmt"

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
	Players 	map[string]PlayerStateInfo
}

type PlayerRanking struct {
	PlayerId 		string
	Ranking 		int
}

type UserRankings struct {
	Rankings 	[]PlayerRanking
}

type DraftPlayerRanking struct {
	// unique player Id will probably just be the team and position such as BUFQB
	PlayerId 		string 	`json:"playerId"`
	// display name for front end
	DisplayName 	string 	`json:"displayName"`
	// team of the player
	Team 			string 	`json:"team"`
	// position of player
	Position 		string 	`json:"position"`
	// address of the user who drafted this player
	OwnerAddress 	string 	`json:"ownerAddress"`
	// number pick that this player was selected.... will default to nil in the database
	PickNum 		int 	`json:"pickNum"`
	// the round which this player was drafted in
	Round 			int 	`json:"round"`
	// rank for the user from mapping rankings onto the players state
	Ranking 		int		`json:"rank"`
}

func CreateRankingObject(rank int, info PlayerStateInfo) DraftPlayerRanking {
	return DraftPlayerRanking{
		PlayerId: info.PlayerId,
		DisplayName: info.DisplayName,
		Team: info.Team,
		Position: info.Position,
		OwnerAddress: info.OwnerAddress,
		PickNum: info.PickNum,
		Round: info.Round,
		Ranking: rank,
	}
}

func GetUserRankings(ownerId string) (*UserRankings, error) {
	var r UserRankings
	err := utils.Db.ReadDocument(fmt.Sprintf("owners/%s/drafts", ownerId), "rankings", &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
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

	res := make(map[string]DraftPlayerRanking)

	for _, rank := range userRankings.Rankings {
		res[rank.PlayerId] = CreateRankingObject(rank.Ranking, state.Players[rank.PlayerId])
	}
	
	return res, nil
}



