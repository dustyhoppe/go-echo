package routing

import (
	"errors"
	"net/http"

	routeErrors "github.com/dustyhoppe/go-echo/routing/errors"
)

type router struct {
	routes []*Route
}

type Router interface {
	MatchRoute(req *http.Request) (*Route, error)
}

func NewRouter(routes []*Route) Router {
	return &router{
		routes: routes,
	}
}

func (r router) MatchRoute(req *http.Request) (*Route, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	var allow []string
	for _, route := range r.routes {
		matches := route.Regex.FindStringSubmatch(req.URL.Path)
		if len(matches) > 0 {
			if req.Method != route.Method {
				allow = append(allow, route.Method)
				continue
			}
			return route, nil
		}
	}

	if len(allow) > 0 {
		return nil, routeErrors.NewMethodNowAllowedError()
	}

	return nil, routeErrors.NewNoMatchingRouteError(req)
}
