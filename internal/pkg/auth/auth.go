// The auth package implements a simple authentication
// mechanism, which allows registration of users, changing
// passwords and login

package auth

import (
	"errors"
	"github.com/boltdb/bolt"
	"html/template"
	"net/http"
	"time"
)

// Errors

var (
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

var (
	ProtectedPaths = []string{"/app/"}
)

var templ = template.Must(template.ParseFiles(
	"web/template/login.html",
	"web/template/register.html"))

// Types

type Session struct {
	Expires time.Time
	User    string
}

type Auth struct {
	sessions map[string]Session
	db       *bolt.DB
}

// Functions

func New() *Auth {
	var a Auth
	db, err := bolt.Open("./data/user.db", 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		panic("Could not open user database")
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		return err
	})
	if err != nil {
		a.db.Close()
		panic("Could not create user bucket")
	}
	a.db = db
	a.sessions = make(map[string]Session)
	return &a
}

func (a *Auth) Close() {
	a.db.Close()
}

// User management

func (a *Auth) AddUser(user, pass string) error {
	return a.db.Update(func(tx *bolt.Tx) error {
		if user == "" || pass == "" {
			return ErrNoUserPass
		}
		b := tx.Bucket([]byte("users"))
		if (b.Get([]byte(user))) == nil {
			hash := hashPassword(user, pass)
			return b.Put([]byte(user), hash[:])
		} else {
			return ErrUserExists
		}
	})
}

func (a *Auth) Login(user, pass string) (string, error) {
	var sess string
	return sess, a.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		if hash := b.Get([]byte(user)); hash == nil {
			return ErrUserNotExists
		} else {
			attempt := hashPassword(user, pass)
			for i, c := range attempt[:] {
				if c != hash[i] {
					return ErrPassword
				}
			}
			// Create session
			sess = sessionString()
			a.sessions[sess] = Session{time.Now().Add(SessionTimeout), user}
			go func() {
				// Remove expired sessions
				time.Sleep(SessionTimeout)
				delete(a.sessions, sess)
			}()
		}
		return nil
	})
}

// Server Auth interface implementation

func (a *Auth) IsValid(sess string) (string, bool) {
	if s := a.sessions[sess]; time.Now().Before(s.Expires) {
		return s.User, true
	}
	return "", false
}

func (a *Auth) Paths() []string {
	return []string{"/login", "/register", "/account"}
}

func (a *Auth) Protect() []string {
	return ProtectedPaths
}

func (a *Auth) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", a.handleLogin)
	mux.HandleFunc("/register", a.handleRegister)
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
