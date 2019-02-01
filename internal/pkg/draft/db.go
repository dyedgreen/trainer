package draft

import (
	"database/sql"
)

func initDb(db *sql.DB) (err error) {
	query := `
	PRAGMA foreign_keys = ON;

	CREATE TABLE IF NOT EXISTS drafts (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		problem INTEGER NOT NULL,
		code TEXT NOT NULL,
		time INTEGER NOT NULL,
		user INTEGER NOT NULL,
		FOREIGN KEY (problem) REFERENCES problems (id),
		FOREIGN KEY (user) REFERENCES users (id)
	);
	`
	_, err = db.Exec(query)
	return
}

func (s *ScratchPad) updateDraft(user int64, draft Draft) (err error) {
	query := `
	DELETE FROM drafts WHERE user = ?;
	INSERT INTO drafts (problem, code, time, user) VALUES (
		?, ?, ?, ?
	);
	`
	_, err = s.db.Exec(query, user, draft.Problem, draft.Code, draft.TimeElapsed, user)
	return
}

func (s *ScratchPad) deleteDraft(user int64) (err error) {
	query := `
	DELETE FROM drafts WHERE user = ?;
	`
	_, err = s.db.Exec(query, user)
	return
}

func (s *ScratchPad) getDraft(user int64) (d Draft, err error) {
	query := `
	SELECT problem, code, time FROM drafts WHERE user = ?;
	`
	row := s.db.QueryRow(query, user)
	err = row.Scan(&d.Problem, &d.Code, &d.TimeElapsed)
	return
}
