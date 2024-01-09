package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func proxyCdn(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "*")
	url := "https://i.kym-cdn.com/" + path

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch image: %s", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to fetch image: Status %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	io.Copy(w, resp.Body)
}

var t = template.Must(template.ParseFiles("ui/article.html"))

func serveArticle(w http.ResponseWriter, r *http.Request) {
	path := r.URL.RequestURI()
	url := "https://knowyourmeme.com" + path

	var result strings.Builder

	c := articleParser(&result)

	c.Visit(url)

	t.ExecuteTemplate(w, "article.html", template.HTML(result.String()))
}
