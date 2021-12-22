package routing

import (
	"fmt"
	"net/http"
	"regexp"
)

type Route struct {
	Method  string
	Regex   *regexp.Regexp
	Handler http.HandlerFunc
}

func NewRoute(method, pattern string, handler http.HandlerFunc) Route {
	return Route{
		Method:  method,
		Regex:   regexp.MustCompile(fmt.Sprintf("^%s$", pattern)),
		Handler: handler,
	}
}
