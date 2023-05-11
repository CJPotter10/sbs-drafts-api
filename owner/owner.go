package owner

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CJPotter10/sbs-drafts-api/models"
	"github.com/go-chi/chi"
)

type OwnerResources struct{}

func (or *OwnerResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/{ownerId}/draftToken/mint/min/{min}/max/{max}", or.CreateTokensInDatabase)
	return r
}

type MintTokensResponse struct {
	Tokens []models.DraftToken `json:"tokens"`
}

func (or *OwnerResources) CreateTokensInDatabase(w http.ResponseWriter, r *http.Request) {
	ownerId := chi.URLParam(r, "ownerId")
	minId := chi.URLParam(r, "min")
	maxId := chi.URLParam(r, "max")
	if ownerId == "" || minId == "" || maxId == "" {
		http.Error(w, "Did not find an ownerId, maxId, or minId in the url path", http.StatusBadRequest)
	}
	min, err := strconv.Atoi(minId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	max, err := strconv.Atoi(maxId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokens := make([]models.DraftToken, 0)
	for i := min; i <= max; i++ {
		tokenId := strconv.Itoa(i)
		token, err := models.MintDraftTokenInDb(tokenId, ownerId)
		if err != nil {
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
