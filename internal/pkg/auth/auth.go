// The auth package implements a simple authentication
// mechanism, which allows registration of users, changing
// passwords and login

package auth

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"
	"strings"
	"time"
)

// Errors

var (
	ErrInit          = errors.New("Could not initialize database")
	ErrUserExists    = errors.New("User already exists")
	ErrPassword      = errors.New("Incorrect password")
	ErrUserNotExists = errors.New("User does not exist")
	ErrNoUserPass    = errors.New("No password or username specified")
)

// Constants and templates

const (
	SessionTimeout   = 24 * time.Hour
	SessionStrLength = 256
)

var templ = template.Must(template.ParseFiles(
	"web/template/login.html",
	"web/template/account.html",
	"web/template/register.html"))

// Types

type Session struct {
	// Hold a session with
	// expiry time and user
	// id
	Expires time.Time
	Username string
	UserId    int64
}

type Auth struct {
	// Authentication provider
	// implementation
	sessions map[string]Session
	db       *sql.DB
}

// Functions

// Create a new auth object
func New(db *sql.DB) *Auth {
	var a Auth
	a.db = db
	a.sessions = make(map[string]Session)
	if err := initDb(a.db); err != nil {
		panic(err.Error())
	}
	return &a
}

// User management

// Add a new user
func (a *Auth) AddUser(user, pass string) error {
	if user == "" || pass == "" {
		return ErrNoUserPass
	} else if a.userExists(user) {
		return ErrUserExists
	}
	var hash, salt string
	salt = randomString(64)
	hash = hashPassword(salt, pass)
	return a.userInsert(user, hash, salt)
}

// Change a users password
func (a *Auth) UpdateUser(user, pass, newPass string) error {
	if user == "" || pass == "" || newPass == "" {
		return ErrNoUserPass
	} else if exists, _, hash, salt := a.userGet(user); !exists {
		return ErrUserNotExists
	} else if hashPassword(salt, pass) != hash {
		return ErrPassword
	}
	// Update user
	var hash, salt string
	salt = randomString(64)
	hash = hashPassword(salt, newPass)
	return a.userUpdate(user, hash, salt)
}

// Create a session for a user
func (a *Auth) Login(user, pass string) (sess string, err error) {
	if exists, userId, hash, salt := a.userGet(user); !exists {
		err = ErrUserNotExists
	} else if hashPassword(salt, pass) != hash {
		err = ErrPassword
	} else {
		sess = randomString(SessionStrLength)
		a.sessions[sess] = Session{time.Now().Add(SessionTimeout), user, userId}
		go func() {
			// Remove expired sessions
			time.Sleep(SessionTimeout)
			delete(a.sessions, sess)
		}()
	}
	return
}

func (a *Auth) GetSession(sess string) (Session, bool) {
	if s := a.sessions[sess]; time.Now().Before(s.Expires) {
		return s, true
	}
	return Session{}, false
}

// Server Auth interface implementation

func (a *Auth) IsValid(sess string) (int64, bool) {
	s, valid := a.GetSession(sess)
	return s.UserId, valid
}

func (a *Auth) Paths() []string {
	return []string{"/login", "/register", "/account"}
}

func (a *Auth) Protect() []string {
	return []string{"/app/", "/api/problem/", "/account"}
}

func (a *Auth) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", a.handleLogin)
	mux.HandleFunc("/register", a.handleRegister)
	mux.HandleFunc("/account", a.handleAccount)
	return mux
}

// Handler functions

func (a *Auth) handleLogin(w http.ResponseWriter, r *http.Request) {
	var message string
	user, pass := r.PostFormValue("username"), r.PostFormValue("password")
	if user != "" {
		sess, err := a.Login(user, pass)
		if err != nil {
			message = err.Error()
		} else {
			var cookie http.Cookie
			cookie.Name = "auth"
			cookie.Value = sess
			cookie.Expires = a.sessions[sess].Expires
			http.SetCookie(w, &cookie)
			r.URL.Path = a.Protect()[0]
			http.Redirect(w, r, r.URL.String(), http.StatusFound)
			return
		}
	}
	templ.ExecuteTemplate(w, "login.html", struct {
		Error    string
		Username string
	}{message, user})
}

func (a *Auth) handleRegister(w http.ResponseWriter, r *http.Request) {
	var message string
	user, pass := r.PostFormValue("username"), r.PostFormValue("password")
	user = strings.Trim(user, " ")
	if user != "" {
		err := a.AddUser(user, pass)
		if err != nil {
			message = err.Error()
		} else {
			r.URL.Path = a.Paths()[0]
			http.Redirect(w, r, r.URL.String(), http.StatusFound)
			return
		}
	}
	templ.ExecuteTemplate(w, "register.html", struct {
		Error    string
		Username string
	}{message, user})
}

func (a *Auth) handleAccount(w http.ResponseWriter, r *http.Request) {
	var message, success string
	pass, newPass := r.PostFormValue("old_pass"), r.PostFormValue("new_pass")
	var sess Session
	for _, c := range r.Cookies() {
		if c.Name == "auth" {
			sess, _ = a.GetSession(c.Value)
			break
		}
	}
	if pass != "" {
		err := a.UpdateUser(sess.Username, pass, newPass)
		if err != nil {
			message = err.Error()
		} else {
			success = "Password updated"
		}
	}
	templ.ExecuteTemplate(w, "account.html", struct {
		Error    string
		Success  string
		Username string
	}{message, success, sess.Username})
}
