package errors

import (
	"fmt"
	"net/http"
)

type NoMatchingRouteError struct {
	request *http.Request
}

func (e *NoMatchingRouteError) Error() string {
	return fmt.Sprintf("no matching route for request '%s %s'", (*e.request).Method, (*e.request).URL.String())
}

func NewNoMatchingRouteError(req *http.Request) *NoMatchingRouteError {
	return &NoMatchingRouteError{
		request: req,
	}
}
