package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	staticServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static", staticServer))

	r.NotFound(pageNotFound)

	r.Get("/cdnproxy/*", proxyCdn)

	r.Get("/", serveMainPage)

	r.Get("/memes/*", serveArticle)

	http.ListenAndServe(":1641", r)
}
