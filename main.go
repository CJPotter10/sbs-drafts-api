package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	draftState "github.com/CJPotter10/sbs-drafts-api/draft-state"
	"github.com/CJPotter10/sbs-drafts-api/leagues"
	"github.com/CJPotter10/sbs-drafts-api/owner"
	"github.com/CJPotter10/sbs-drafts-api/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {

	port := "8080"

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	utils.NewDatabaseClient()

	fmt.Printf("Starting up on http://localhost:%s\n", port)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	dr := &draftState.DraftResources{}
	r.Mount("/drafts", dr.Routes())

	lr := &leagues.LeagueResources{}
	r.Mount("/league", lr.Routes())

	or := &owner.OwnerResources{}
	r.Mount("/owner", or.Routes())

	log.Fatal(http.ListenAndServe(":"+port, r))
}
