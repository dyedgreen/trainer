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
	"sync"
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
	TicketKeyLength  = 16
)

var templ = template.Must(template.ParseFiles(
	"web/template/login.html",
	"web/template/account.html",
	"web/template/register.html",
	"web/template/tickets.html"))

// Types

type Session struct {
	// Hold a session with
	// expiry time and user
	// id
	Expires  time.Time
	Username string
	UserId   int64
}

type Auth struct {
	// Authentication provider
	// implementation
	sessions map[string]Session
	mutex    sync.RWMutex
	db       *sql.DB
	protect []string
}

// Functions

// Create a new auth object
func New(db *sql.DB, protect []string) *Auth {
	var a Auth
	a.db = db
	a.sessions = make(map[string]Session)
	if err := initDb(a.db); err != nil {
		panic(err.Error())
	}
	a.protect = append(protect, "/account", "/tickets")
	return &a
}

// User management

// Add a new user
func (a *Auth) AddUser(user, ticket, pass string) error {
	if user == "" || pass == "" {
		return ErrNoUserPass
	} else if a.userExists(user) {
		return ErrUserExists
	} else if err := a.ticketUse(ticket); err != nil {
		// The first person to register is admin
		if count, errCount := a.userCount(); errCount != nil || count > 0 {
			return err
		}
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
		a.mutex.Lock()
		defer a.mutex.Unlock()
		a.sessions[sess] = Session{time.Now().Add(SessionTimeout), user, userId}
		go func() {
			// Remove expired sessions
			time.Sleep(SessionTimeout)
			a.mutex.Lock()
			defer a.mutex.Unlock()
			delete(a.sessions, sess)
		}()
	}
	return
}

func (a *Auth) GetSession(sess string) (Session, bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
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
	return []string{"/login", "/register", "/account", "/logout", "/tickets"}
}

func (a *Auth) Protect() []string {
	return a.protect
}

func (a *Auth) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", a.handleLogin)
	mux.HandleFunc("/register", a.handleRegister)
	mux.HandleFunc("/account", a.handleAccount)
	mux.HandleFunc("/logout", a.handleLogout)
	mux.HandleFunc("/tickets", a.handleTickets)
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
			a.mutex.RLock()
			defer a.mutex.RUnlock()
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
	user, ticket, pass := r.PostFormValue("username"), r.PostFormValue("ticket"), r.PostFormValue("password")
	user = strings.Trim(user, " ")
	if user != "" {
		err := a.AddUser(user, ticket, pass)
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
		Ticket   string
	}{message, user, ticket})
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

func (a *Auth) handleLogout(w http.ResponseWriter, r *http.Request) {
	var cookie http.Cookie
	cookie.Name = "auth"
	cookie.Value = ""
	cookie.Expires = time.Now()
	http.SetCookie(w, &cookie)
	r.URL.Path = "/"
	http.Redirect(w, r, r.URL.String(), http.StatusFound)
}

func (a *Auth) handleTickets(w http.ResponseWriter, r *http.Request) {
	var err error
	var sess Session
	for _, c := range r.Cookies() {
		if c.Name == "auth" {
			sess, _ = a.GetSession(c.Value)
			break
		}
	}
	if sess.UserId != 1 {
		// Only the first user is an admin
		templ.ExecuteTemplate(w, "tickets.html", struct {
			Error   string
			Tickets []string
		}{"Forbidden", nil})
		return
	}
	if r.PostFormValue("new") == "yes" {
		err = a.ticketInsert(randomString(TicketKeyLength))
	}
	var tickets []string
	var message string
	if err == nil {
		tickets, err = a.ticketList()
	}
	if err != nil {
		message = err.Error()
	}
	templ.ExecuteTemplate(w, "tickets.html", struct {
		Error   string
		Tickets []string
	}{message, tickets})
}
