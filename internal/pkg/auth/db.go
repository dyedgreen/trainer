// Auth helpers related to db
// interaction

package auth

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrInvalidKey = errors.New("Invalid key")
	ErrNoTicket   = errors.New("Ticket does not exist")
)

func initDb(db *sql.DB) (err error) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(64) NOT NULL UNIQUE,
		password VARCHAR(256),
		salt VARCHAR(64)
	);

	CREATE TABLE IF NOT EXISTS tickets (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		key VARCHAR(16) NOT NULL UNIQUE,
		used INTEGER NOT NULL,
		date INTEGER
	);
	`
	_, err = db.Exec(query)
	return
}

func (a *Auth) userExists(username string) bool {
	query := `
	SELECT COUNT(*) FROM users WHERE username = ?
	`
	row := a.db.QueryRow(query, username)
	var count int
	if err := row.Scan(&count); err != nil {
		panic(err.Error())
	}
	return count == 1
}

func (a *Auth) userGet(username string) (exists bool, id int64, password, salt string) {
	query := `
	SELECT id, password, salt FROM users WHERE username = ?
	`
	row := a.db.QueryRow(query, username)
	exists = row.Scan(&id, &password, &salt) == nil
	return
}

func (a *Auth) userInsert(username, password, salt string) (err error) {
	query := `
	INSERT INTO users (username, password, salt) VALUES (?, ?, ?)
	`
	_, err = a.db.Exec(query, username, password, salt)
	return
}

func (a *Auth) userUpdate(username, password, salt string) (err error) {
	query := `
	UPDATE users SET password = ?, salt = ? WHERE username = ?
	`
	_, err = a.db.Exec(query, password, salt, username)
	return
}

func (a *Auth) userCount() (count int, err error) {
	query := `
	SELECT COUNT(*) FROM users
	`
	row := a.db.QueryRow(query)
	err = row.Scan(&count)
	return
}

func (a *Auth) ticketInsert(key string) (err error) {
	if len(key) != TicketKeyLength {
		return ErrInvalidKey
	}
	query := `
	INSERT INTO tickets (key, used, date) VALUES (?, 0, 0)
	`
	_, err = a.db.Exec(query, key)
	return
}

func (a *Auth) ticketUse(key string) (err error) {
	query := `
	SELECT COUNT(*) FROM tickets WHERE key = ? AND used = 0
	`
	row := a.db.QueryRow(query, key)
	var count int
	if err = row.Scan(&count); err != nil || count == 0 {
		return ErrNoTicket
	}
	// Update ticket to used
	query = `
	UPDATE tickets SET used = 1, date = ? WHERE key = ?
	`
	_, err = a.db.Exec(query, time.Now().Unix(), key)
	return
}

func (a *Auth) ticketList() ([]string, error) {
	query := `
	SELECT key FROM tickets WHERE used = 0
	`
	var list []string
	rows, err := a.db.Query(query)
	if err != nil {
		return list, err
	}
	var key string
	for rows.Next() {
		err = rows.Scan(&key)
		if err != nil {
			return list, err
		}
		list = append(list, key)
	}
	return list, nil
}
