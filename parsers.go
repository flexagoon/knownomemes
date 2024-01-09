package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func kymParser(goquerySelector string, f colly.HTMLCallback) *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains("knowyourmeme.com"),
	)

	c.OnHTML(goquerySelector, f)

	return c
}

func articleParser(result *strings.Builder) *colly.Collector {
	return kymParser("#maru article.entry", func(h *colly.HTMLElement) {
		result.WriteString("<hgroup>\n")
		title := h.ChildText("h1")
		result.WriteString(wrapWithTag(title, "h1"))
		coverURL := h.ChildAttr("header .photo img", "src")
		cover := proxyImage(coverURL)
		result.WriteString(cover)
		result.WriteString("</hgroup>\n")
		h.ForEach("#entry_body section.bodycopy :not(#search-interest)", func(i int, h *colly.HTMLElement) {
			if h.Name == "h2" {
				result.WriteString(wrapWithTag(h.Text, "h2"))
			} else if h.Name == "p" {
				if !h.DOM.Parent().HasClass("references") {
					text := linkify(h)
					if text != "" {
						result.WriteString(wrapWithTag(text, "p"))
					}
				} else {
					footnote := parseFootnote(h)
					result.WriteString(footnote)
				}
			} else if h.Name == "lite-youtube" {
				id := h.Attr("videoid")
				embed := fmt.Sprintf("<iframe src='https://yewtu.be/embed/%s?autoplay=0&local=true' frameborder='0'></iframe>\n", id)
				result.WriteString(embed)
			} else if h.Name == "lite-tiktok" {
				id := h.Attr("videoid")
				tiktokUrl := fmt.Sprintf("https://tiktok.com/@knm/video/%s", id)
				proxitokUrl := fmt.Sprintf("https://proxitok.belloworld.it/@knm/video/%s", id)
				text := fmt.Sprintf("<a href='%s'>♪ TikTok Video</a> (<a href='%s'>Mirror</a>)", tiktokUrl, proxitokUrl)
				result.WriteString(wrapWithTag(text, "p"))
			} else if h.Name == "img" {
				src := h.Attr("data-src")
				img := proxyImage(src)
				result.WriteString(img)
			}
		})
	})
}
