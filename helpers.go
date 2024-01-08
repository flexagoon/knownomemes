package main

import (
	"fmt"
	"strings"
)

func wrapWithTag(text string, tag string) string {
	return fmt.Sprintf("<%[1]s>%[2]s</%[1]s>\n", tag, text)
}

func proxyImage(src string) string {
	proxy := strings.Replace(src, "https://i.kym-cdn.com", "/cdnproxy", 1)
	img := fmt.Sprintf("<img src='%s' width='45%%'>\n", proxy)
	return img
}
