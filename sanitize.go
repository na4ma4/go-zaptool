package zaptool

import (
	"html"
	"strings"
)

var (
	//nolint:gochecknoglobals // same as html.htmlEscaper.
	userAgentEscaper = strings.NewReplacer(
		"\n", "\\n",
		"+", " ",
		"\000", "\\0",
	)
	//nolint:gochecknoglobals // same as html.htmlEscaper.
	uriEscaper = strings.NewReplacer(
		"\n", "\\n",
		"\000", "\\0",
	)
)

func sanitizeUserAgent(ua string) string {
	return userAgentEscaper.Replace(ua)
}

func sanitizeURI(url string) string {
	return uriEscaper.Replace(url)
}

func sanitizeUsername(username string) string {
	return html.EscapeString(username)
}
