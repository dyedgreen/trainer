// Package server provides a file server,
// which can be registered with an auth
// provider and optional ApiHandlers

package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
)

const (
	staticRoot   = "./web/static"
	templateRoot = "./web/template"
)

// Interfaces

type Auth interface {
	// Authentication provider interface
	// the authentication provider manages
	// sessions. Sessions are represented by
	// strings and stored in the cookie 'auth'
	//
	// This interface can register a handler
	// for a set of paths to handle login /
	// register etc. pages.
	//
	// The first path is expected to be the
	// login page. It is the implementations
	// responsibility to store the session string
	// in the 'auth' cookie.
	IsValid(session string) (string, bool)
	Paths() []string
	Protect() []string
	Handler() http.Handler
}

type Api interface {
	// Api methods can be registered. They
	// handle /api{/.Path()} and responses are
	// returned in JSON format
	//
	// The handling function is passed the user
	// identifier (or empty string if not logged
	// in)
	Path() string
	Call(r *http.Request, user string) (interface{}, error)
}

// Page Templates

var templ = template.Must(template.ParseFiles(templateRoot + "/error.html"))

// Server

type Server struct {
	server http.Server
	mux    *http.ServeMux
	auth   Auth
}

func New(addr string) *Server {
	var s Server
	s.server.Addr = addr
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/", s.handleFile)
	s.server.Handler = s.mux
	return &s
}

// Handler functions

func (s *Server) handleError(status int, w http.ResponseWriter) {
	w.WriteHeader(status)
	templ.ExecuteTemplate(w, "error.html", struct {
		Code int
		Text string
	}{status, http.StatusText(status)})
}

func (s *Server) handleFile(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat(staticRoot + r.URL.Path); err != nil {
		s.handleError(http.StatusNotFound, w)
		return
	}
	http.ServeFile(w, r, staticRoot+r.URL.Path)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if s.auth == nil {
		s.handleError(http.StatusInternalServerError, w)
		return
	}
	var cookie *http.Cookie
	var hasSess bool
	for _, c := range r.Cookies() {
		if c.Name == "auth" {
			cookie = c
			hasSess = true
		}
	}
	if hasSess {
		if _, valid := s.auth.IsValid(cookie.Value); valid {
			s.handleFile(w, r)
			return
		}
	}
	if len(s.auth.Paths()) > 0 {
		r.URL.Path = s.auth.Paths()[0]
		http.Redirect(w, r, r.URL.String(), http.StatusTemporaryRedirect)
	} else {
		s.handleError(http.StatusUnauthorized, w)
	}
}

// Functions used to mutate the server object

func (s *Server) RegisterAuth(auth Auth) {
	if s.auth != nil {
		panic("Authentication provider already registered")
	}
	s.auth = auth
	for _, path := range s.auth.Paths() {
		s.mux.Handle(path, s.auth.Handler())
	}
	for _, path := range s.auth.Protect() {
		s.mux.HandleFunc(path, s.handleLogin)
	}
}

func (s *Server) RegisterApi(api Api) {
	s.mux.HandleFunc("/api"+api.Path(), func(w http.ResponseWriter, r *http.Request) {
		var user string
		if cookie, err := r.Cookie("auth"); err == nil {
			user, _ = s.auth.IsValid(cookie.Value)
		}

		var response, result interface{}
		result, err := api.Call(r, user)

		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response = struct {
				Error string      `json:"error"`
				Value interface{} `json:"value"`
			}{err.Error(), result}
		} else {
			response = struct {
				Error bool        `json:"error"`
				Value interface{} `json:"value"`
			}{false, result}
		}
		json.NewEncoder(w).Encode(response)
	})
}

// Functions used to run the server

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}
