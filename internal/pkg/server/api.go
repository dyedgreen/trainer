// Functions that help when dealing with
// api registration

package server

import (
	"errors"
	"net/http"
)

var (
	ErrApiNotFound = errors.New("Endpoint not found")
)

type ApiPlaceholder string

func (a ApiPlaceholder) Path() string {
	return string(a)
}

func (a ApiPlaceholder) Call(r *http.Request, user int64) (interface{}, error) {
	return nil, ErrApiNotFound
}

type apiFunction struct {
	f    func(*http.Request, int64) (interface{}, error)
	path string
}

func (a *apiFunction) Path() string {
	return a.path
}

func (a *apiFunction) Call(r *http.Request, user int64) (interface{}, error) {
	return a.f(r, user)
}
