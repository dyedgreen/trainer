// Auth helpers related to db
// interaction

package auth

import (
	"database/sql"
)

func initDb(db *sql.DB) (err error) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(64) NOT NULL UNIQUE,
		password VARCHAR(256),
		salt VARCHAR(64)
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
