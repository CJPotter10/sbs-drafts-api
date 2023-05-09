package owner

import (
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
	for i := min; i < max; i++ {
		tokenId := strconv.Itoa(i)
		err := models.MintDraftTokenInDb(tokenId, ownerId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
