package draftState

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CJPotter10/sbs-drafts-api/models"
	"github.com/go-chi/chi"
)

type DraftResources struct{}

func (dr *DraftResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/{draftId}/playerState/{ownerId}", dr.getPlayersMapWithRankings)
	r.Get("/{draftId}/state/info", dr.getDraftInfoById)
	r.Get("/{draftId}/state/summary", dr.getDraftSummaryById)
	r.Get("/{draftId}/state/connectionList", dr.getDraftConnectionList)
	r.Get("/{draftId}/state/rosters", dr.getRostersMapForDraft)
	return r
}

// will need to add the stats and analysis that needs to be shown to this route when we have that data
func (dr *DraftResources) getPlayersMapWithRankings(w http.ResponseWriter, r *http.Request) {
	ownerId := chi.URLParam(r, "ownerId")
	draftId := chi.URLParam(r, "draftId")
	if ownerId == "" {
		http.Error(w, "Did not find an ownerId in this request so we are returning", http.StatusBadRequest)
		fmt.Println("Did not find an ownerId in this request so we are returning")
		return
	}
	if draftId == "" {
		http.Error(w, "Did not find an draftId in this request so we are returning", http.StatusBadRequest)
		fmt.Println("Did not find an draftId in this request so we are returning")
		return
	}

	res, err := models.ReturnPlayerStateWithRankings(ownerId, draftId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("first entry of the array: ", res[0])

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

// needs to return draft info such as what pick we are on

func (dr *DraftResources) getDraftInfoById(w http.ResponseWriter, r *http.Request) {
	draftId := chi.URLParam(r, "draftId")
	if draftId == "" {
		http.Error(w, "No draft Id was found in the URL", 400)
		return
	}

	info, err := models.ReturnDraftInfoForDraft(draftId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(info)
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

// returns the draft summary
func (dr *DraftResources) getDraftSummaryById(w http.ResponseWriter, r *http.Request) {
	draftId := chi.URLParam(r, "draftId")
	if draftId == "" {
		http.Error(w, "No draft Id was found in the URL", 400)
		return
	}

	sum, err := models.ReturnDraftSummaryForDraft(draftId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(sum)
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

// returns map of connection list
func (dr *DraftResources) getDraftConnectionList(w http.ResponseWriter, r *http.Request) {
	draftId := chi.URLParam(r, "draftId")
	if draftId == "" {
		http.Error(w, "No draft Id was found in the URL", 400)
		return
	}

	cl, err := models.ReturnConnectionListForDraft(draftId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(cl)
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

// returns rosters for all users in the draft
func (dr *DraftResources) getRostersMapForDraft(w http.ResponseWriter, r *http.Request) {
	draftId := chi.URLParam(r, "draftId")
	if draftId == "" {
		http.Error(w, "No draft Id was found in the URL", 400)
		return
	}

	rs, err := models.ReturnRostersForDraft(draftId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(rs)
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
