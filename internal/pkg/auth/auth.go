// The auth package implements a simple authentication
// mechanism, which allows registration of users, changing
// passwords and login

package auth

import (
	"net/http"
	"time"
	"html/template"
)

// Constants and templates

const (
	SessionTimeout = 86400 * time.Second
)

var templ = template.Must(template.ParseFiles(templateRoot + "/error.html"))

// Types

type Session struct {
	Valid time.Time
	User  string
}

type Auth struct {
	sessions map[string]Session
}

// Functions

func New() *Auth {
	return new(Auth)
}

func (a *Auth) IsValid(sess string) (string, bool) {
	if s := a.sessions[sess]; s.Valid.Unix() > time.Now().Unix() {
		return s.User, true
	}
	return "", false
}

func (a *Auth) Paths() []string {
	return []string{"login", "register", "account"}
}

func (a *Auth) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", a.handleLogin)
	return mux
}

// Handler functions

func (a *Auth) handleLogin(w http.ResponseWriter, r *http.Request) {

}
