package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func wrapWithTag(text string, tag string) string {
	return fmt.Sprintf("<%[1]s>%[2]s</%[1]s>\n", tag, text)
}

func proxyImage(src string) string {
	proxy := strings.Replace(src, "https://i.kym-cdn.com", "/cdnproxy", 1)
	img := fmt.Sprintf("<img src='%s' width='45%%'>\n", proxy)
	return img
}

func linkify(h *colly.HTMLElement) string {
	text := h.DOM.Contents().Map(func(i int, h *goquery.Selection) string {
		nodeType := goquery.NodeName(h)
		if nodeType == "a" {
			href, _ := h.Attr("href")
			href = strings.Replace(href, "https://knowyourmeme.com", "", -1)
			text := h.Text()
			return fmt.Sprintf("<a href='%s'>%s</a>", href, text)
		} else if nodeType == "sup" {
			a := h.ChildrenFiltered("a")
			href, _ := a.Attr("href")
			fnId := href[3:] // Strip "#fn" from href
			text := a.Text()
			return fmt.Sprintf("<sup id='fns%[1]s'><a href='#fn%[1]s'>%s</a></sup>", fnId, text)
		}
		return h.Text()
	})
	return strings.Join(text, "")
}

func parseFootnote(h *colly.HTMLElement) string {
	fn := h.ChildText("sup a")
	fnId := fn[1 : len(fn)-1]
	sourceSite := strings.SplitAfterN(h.ChildText(".footnote-text"), " – ", 2)[0]
	sourceUrl := h.ChildAttr(".footnote-text a", "href")
	sourceText := h.ChildText(".footnote-text a")
	return fmt.Sprintf(
		"<p id='fn%[1]s'><a href='#fns%[1]s'>[%[1]s]</a> %s<a href='%s'>%s</a></p>\n",
		fnId,
		sourceSite,
		sourceUrl,
		sourceText,
	)
}
