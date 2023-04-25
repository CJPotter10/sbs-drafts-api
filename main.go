package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	draftState "github.com/CJPotter10/sbs-drafts-api/draft-state"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	port := "8888"

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	fmt.Printf("Starting up on http://localhost:%s\n", port)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func( w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	dr := &draftState.DraftResources{}
	r.Mount("/drafts", dr.Routes())

	log.Fatal(http.ListenAndServe(":"+port, r))
}