package leagues

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CJPotter10/sbs-drafts-api/models"
	"github.com/go-chi/chi"
)

type LeagueResources struct{}

func (lr *LeagueResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/{draftType}/owner/{ownerId}", lr.joinDraftLeagues)
	r.Post("/{draftId}/actions/leave", lr.RemoveUserFromDraft)
	r.Get("/{draftId}/cards/{tokenId}", lr.ReturnDraftToken)
	return r
}

type JoinLeagueRequestBody struct {
	NumLeaguesToJoin int `json:"numLeaguesToJoin"`
}

// route to join draft league
func (lr *LeagueResources) joinDraftLeagues(w http.ResponseWriter, r *http.Request) {
	ownerId := chi.URLParam(r, "ownerId")
	draftType := chi.URLParam(r, "draftType")
	if ownerId == "" || draftType == "" {
		http.Error(w, "Did not find an ownerid in this request so we are returning", http.StatusInternalServerError)
		return
	}

	var req JoinLeagueRequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("Error in decoding the request body for joining leagues")
		http.Error(w, "Could not decode the request body into the correct data type so we are returning", http.StatusInternalServerError)
		return
	}

	cards, err := models.JoinLeagues(ownerId, req.NumLeaguesToJoin, draftType)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(cards)
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

// route to leave draft league (make sure that the draft has not started or alread happened)
type LeaveRequest struct {
	OwnerId string `json:"ownerId"`
	TokenId string `json:"tokenId"`
}

func (lr *LeagueResources) RemoveUserFromDraft(w http.ResponseWriter, r *http.Request) {
	draftId := chi.URLParam(r, "draftId")

	var req LeaveRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("Error in decoding the request body for leaving league")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = models.RemoveUserFromDraft(req.TokenId, req.OwnerId, draftId)
	if err != nil {
		fmt.Println("Error in decoding the request body for leaving league")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// route to return roster/card for user to be used after the draft
func (lr *LeagueResources) ReturnDraftToken(w http.ResponseWriter, r *http.Request) {
	draftId := chi.URLParam(r, "draftId")
	tokenId := chi.URLParam(r, "tokenId")
	if draftId == "" || tokenId == "" {
		http.Error(w, "The draftId or TokenID that was passed in were empty", http.StatusBadRequest)
		return
	}

	var t models.DraftToken
	err := t.GetDraftTokenFromDraftById(tokenId, draftId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// route to return the leaderboard for draft leagues

// route to return leaderboard for all of the draft leagues top scores for a gameweek
