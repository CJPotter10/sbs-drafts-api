package owner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CJPotter10/sbs-drafts-api/models"
	"github.com/CJPotter10/sbs-drafts-api/utils"
	"github.com/go-chi/chi"
)

type OwnerResources struct{}

func (or *OwnerResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/{ownerId}/draftToken/mint/min/{min}/max/{max}", or.CreateTokensInDatabase)
	r.Post("/{ownerId}/drafts/state/rankings", or.UpdateUserRankings)
	r.Get("/{ownerId}/draftToken/all", or.ReturnTokensOwnedByUser)
	r.Get("/{ownerId}/rankings/get", or.ReturnUserRankings)
	return r
}

type MintTokensResponse struct {
	Tokens []models.DraftToken `json:"tokens"`
}

func (or *OwnerResources) CreateTokensInDatabase(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside of request")
	ownerId := chi.URLParam(r, "ownerId")
	minId := chi.URLParam(r, "min")
	maxId := chi.URLParam(r, "max")
	if ownerId == "" || minId == "" || maxId == "" {
		fmt.Println("no urls were found")
		http.Error(w, "Did not find an ownerId, maxId, or minId in the url path", http.StatusBadRequest)
	}
	min, err := strconv.Atoi(minId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	max, err := strconv.Atoi(maxId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokens := make([]models.DraftToken, 0)
	for i := min; i <= max; i++ {
		tokenId := strconv.Itoa(i)
		token, err := models.MintDraftTokenInDb(tokenId, ownerId)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tokens = append(tokens, *token)
	}

	res := &MintTokensResponse{
		Tokens: tokens,
	}

	data, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (or *OwnerResources) ReturnTokensOwnedByUser(w http.ResponseWriter, r *http.Request) {
	ownerId := chi.URLParam(r, "ownerId")
	if ownerId == "" {
		http.Error(w, "Did not find an ownerId in the url path", http.StatusInternalServerError)
		return
	}

	res, err := models.ReturnAllDraftTokensForOwner(ownerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (or *OwnerResources) UpdateUserRankings(w http.ResponseWriter, r *http.Request) {
	ownerId := chi.URLParam(r, "ownerId")
	if ownerId == "" {
		http.Error(w, "Did not find an ownerId in the url path", http.StatusInternalServerError)
		return
	}

	var newRankings models.UserRankings
	err := json.NewDecoder(r.Body).Decode(&newRankings)
	if err != nil {
		fmt.Println("Error in decoding the request body for updating this users rankings")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.Db.CreateOrUpdateDocument(fmt.Sprintf("owners/%s/drafts", ownerId), "rankings", newRankings)
	if err != nil {
		fmt.Println("error in updating the owners rankings in the db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(newRankings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type GetRankingsResponse struct {
	PlayerId string             `json:"playerId"`
	Rank     int                `json:"rank"`
	Score    int                `json:"score"`
	Stats    models.StatsObject `json:"stats"`
}

func (or *OwnerResources) ReturnUserRankings(w http.ResponseWriter, r *http.Request) {
	ownerId := chi.URLParam(r, "ownerId")
	if ownerId == "" {
		http.Error(w, "Did not find an ownerId in the url path", http.StatusInternalServerError)
		return
	}

	res, err := models.GetUserRankings(ownerId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var stats models.StatsMap
	err = utils.Db.ReadDocument("playerStats2023", "playerMap", &stats)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]GetRankingsResponse, 0)

	for i := 0; i < len(res.Rankings); i++ {
		obj := GetRankingsResponse{
			PlayerId: res.Rankings[i].PlayerId,
			Rank:     res.Rankings[i].Ranking,
			Score:    res.Rankings[i].Score,
			Stats:    stats.PlayerStats[res.Rankings[i].PlayerId],
		}
		response = append(response, obj)
	}

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
