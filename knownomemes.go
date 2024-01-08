package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/cdnproxy/*", proxyCdn)

	r.Get("/memes/*", serveArticle)

	http.ListenAndServe(":1641", r)
}
