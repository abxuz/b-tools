package bhttp

import (
	"net/http"
	"strings"
)

type Object = map[string]any
type Array = []Object

func Scheme(req *http.Request) string {
	scheme := strings.ToLower(req.Header.Get("X-Forwarded-Proto"))
	switch scheme {
	case "http", "https":
		return scheme
	}
	if req.TLS == nil {
		return "http"
	}
	return "https"
}
