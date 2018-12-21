// Problem helpers related
// to db interaction

package problem

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrEmpty            = errors.New("Values may not be empty")
	ErrProblemNotExists = errors.New("Problem does not exist")
)

func initDb(db *sql.DB) (err error) {
	query := `
	PRAGMA foreign_keys = ON;

	CREATE TABLE IF NOT EXISTS problems (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title VARCHAR(64),
		question TEXT NOT NULL,
		solution TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS schedule (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		problem INTEGER NOT NULL,
		user INTEGER NOT NULL,
		due INTEGER NOT NULL,
		FOREIGN KEY (problem) REFERENCES problems (id),
		FOREIGN KEY (user) REFERENCES users (id)
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		problem INTEGER NOT NULL,
		user INTEGER NOT NULL,
		date INTEGER NOT NULL,
		code TEXT NOT NULL,
		time INTEGER NOT NULL,
		solved INTEGER NOT NULL,
		FOREIGN KEY (problem) REFERENCES problems (id),
		FOREIGN KEY (user) REFERENCES users (id)
	);
	`
	_, err = db.Exec(query)
	return
}

func (b *Box) createProblem(p Problem) (Problem, error) {
	if p.Title == "" || p.Question == "" || p.Solution == "" {
		return p, ErrEmpty
	}
	query := `
	INSERT INTO problems (title, question, solution) VALUES (?, ?, ?);`
	if res, err := b.db.Exec(query, p.Title, p.Question, p.Solution); err != nil {
		return p, err
	} else {
		p.Id, _ = res.LastInsertId()
		return p, err
	}
}

func (b *Box) updateProblem(p Problem) error {
	if p.Title == "" || p.Question == "" || p.Solution == "" {
		return ErrEmpty
	}
	query := `UPDATE problems SET title = ?, question = ?, solution = ? WHERE id = ?;`
	if res, err := b.db.Exec(query, p.Title, p.Question, p.Solution, p.Id, p.Id); err != nil {
		return err
	} else if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrProblemNotExists
	}
	return nil
}

func (b *Box) storeSession(s Session) (err error) {
	if s.Time < 1 || s.Code == "" {
		err = ErrEmpty
		return
	}
	query := `
	INSERT INTO sessions (problem, user, date, code, time, solved) values (
		?,
		(SELECT id FROM users WHERE username = ? LIMIT 1),
		?,
		?,
		?,
		?
	);
	`
	_, err = b.db.Exec(query, s.Problem, s.User, s.Date, s.Code, s.Time, s.Solved)
	return
}

func (b *Box) numSuccessfulAttempts(id int64, user string) (num int) {
	query := `
	SELECT COUNT(*) FROM sessions WHERE solved = 1 AND problem = ? AND user = (
		SELECT id FROM users WHERE username = ? LIMIT 1
	) AND date > IFNULL(
		(SELECT date FROM sessions WHERE solved != 1 AND problem = ? AND user = (
			SELECT id FROM users WHERE username = ? LIMIT 1
		) ORDER BY date DESC LIMIT 1),
		0
	);
	`
	// On error, num is 0
	b.db.QueryRow(query, id, user, id, user).Scan(&num)
	return
}

func (b *Box) scheduleProblem(id int64, user string, due int64) (err error) {
	query := `
	DELETE FROM schedule WHERE problem = ? AND user = (
		SELECT id FROM users WHERE username = ? LIMIT 1
	);

	INSERT INTO schedule (problem, user, due) VALUES (
		?,
		(SELECT id FROM users WHERE username = ? LIMIT 1),
		?
	);
	`
	_, err = b.db.Exec(query, id, user, id, user, due)
	return
}

func (b *Box) nextScheduledProblem(user string) (p Problem, err error) {
	query := `
	SELECT id, title, question, solution FROM problems WHERE id = (
		SELECT problem FROM schedule WHERE due <= ? AND user = (
			SELECT id FROM users WHERE username = ? LIMIT 1
		) ORDER BY due ASC LIMIT 1
	);
	`
	row := b.db.QueryRow(query, time.Now().Unix(), user)
	err = row.Scan(&p.Id, &p.Title, &p.Question, &p.Solution)
	return
}

func (b *Box) notScheduledProblem(user string) (p Problem, err error) {
	query := `
	SELECT id, title, question, solution FROM problems WHERE NOT id IN (
		SELECT problem FROM schedule WHERE user = (
			SELECT id FROM users WHERE username = ?
		)
	) ORDER BY RANDOM() LIMIT 1;
	`
	row := b.db.QueryRow(query, user)
	err = row.Scan(&p.Id, &p.Title, &p.Question, &p.Solution)
	return
}
