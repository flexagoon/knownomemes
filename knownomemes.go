package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gocolly/colly"
)

func main() {
	r := chi.NewRouter()

	t := template.Must(template.ParseFiles("article.html"))

	r.Get("/cdnproxy/*", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.RequestURI()
		url := "https://i.kym-cdn.com" + path

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
	})

	r.Get("/memes/*", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.RequestURI()
		url := "https://knowyourmeme.com" + path

		var result strings.Builder

		c := colly.NewCollector(
			colly.AllowedDomains("knowyourmeme.com"),
		)

		c.OnHTML("#maru article.entry", func(h *colly.HTMLElement) {
			result.WriteString("<hgroup>\n")
			title := h.ChildText("h1")
			result.WriteString(wrapWithTag(title, "h1"))
			coverURL := h.ChildAttr("header .photo img", "src")
			cover := fmt.Sprintf("<img src='%s' width='45%%'>\n", coverURL)
			result.WriteString(cover)
			result.WriteString("</hgroup>\n")
			h.ForEach("#entry_body section.bodycopy :not(#search-interest)", func(i int, h *colly.HTMLElement) {
				if h.Name == "p" || h.Name == "h2" {
					result.WriteString(wrapWithTag(h.Text, h.Name))
					if h.DOM.Parent().HasClass("references") {
						fmt.Println(h.Text)
					}
				} else if h.Name == "lite-youtube" {
					id := h.Attr("videoid")
					embed := fmt.Sprintf("<iframe src='https://yewtu.be/embed/%s?autoplay=0' frameborder='0'></iframe>\n", id)
					result.WriteString(embed)
				} else if h.Name == "img" {
					src := h.Attr("data-src")
					img := fmt.Sprintf("<img src='%s' width='45%%'>\n", src)
					result.WriteString(img)
				}
			})
		})

		c.Visit(url)

		t.ExecuteTemplate(w, "article.html", template.HTML(result.String()))
	})

	http.ListenAndServe(":1641", r)
}

func wrapWithTag(text string, tag string) string {
	return fmt.Sprintf("<%[1]s>%[2]s</%[1]s>\n", tag, text)
}
