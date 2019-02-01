// Allows users to persist a draft for
// the current problem

package draft

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
)

var (
	ErrDraftNotExists = errors.New("User has now draft")
	ErrNeedProblem    = errors.New("Problem is not specified")
	ErrNeedTime       = errors.New("Times elapsed is not specified")
)

type Draft struct {
	Problem     int64  `json:"problem"` // Problem id worked on
	Code        string `json:"code"`    // Code draft
	TimeElapsed int64  `json:"time"`    // Time elapsed in seconds
}

type ScratchPad struct {
	db *sql.DB
}

func NewScratchPad(db *sql.DB) *ScratchPad {
	if err := initDb(db); err != nil {
		panic(err.Error())
	}
	var pad ScratchPad
	pad.db = db
	return &pad
}

// Functions exposed to api

func (s *ScratchPad) DraftUpdate(r *http.Request, user int64) (interface{}, error) {
	var draft Draft
	var err error
	if draft.Problem, err = strconv.ParseInt(r.FormValue("problem"), 10, 64); err != nil {
		return nil, ErrNeedProblem
	}
	if draft.TimeElapsed, err = strconv.ParseInt(r.FormValue("time"), 10, 64); err != nil {
		return nil, ErrNeedTime
	}
	draft.Code = r.FormValue("code")
	return nil, s.updateDraft(user, draft)
}

func (s *ScratchPad) DraftGet(r *http.Request, user int64) (interface{}, error) {
	if draft, err := s.getDraft(user); err == nil {
		return draft, nil
	}
	return nil, ErrDraftNotExists
}
