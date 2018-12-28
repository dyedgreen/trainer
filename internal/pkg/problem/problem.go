// The problem package manages problems for
// users

package problem

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Problem struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Question string `json:"question"`
	Solution string `json:"solution"`
}

type Session struct {
	Id      int64
	Problem int64
	User    int64
	Date    int64
	Code    string
	Time    int64
	Solved  bool
}

type Box struct {
	// Contains problems
	db *sql.DB
}

func NewBox(db *sql.DB) *Box {
	if err := initDb(db); err != nil {
		panic(err.Error())
	}
	var b Box
	b.db = db
	return &b
}

// Implement server api functions

func (b *Box) ProblemUpdate(r *http.Request, user int64) (interface{}, error) {
	// Update (or create) problem and return problem id
	var id int64
	var err error
	if id, err = strconv.ParseInt(r.FormValue("id"), 10, 64); err != nil {
		return nil, err
	}
	problem := Problem{id, r.FormValue("title"), r.FormValue("question"), r.FormValue("solution")}
	problem.Title = strings.Trim(problem.Title, " \n")
	problem.Question = strings.Trim(problem.Question, " \n")
	problem.Solution = strings.Trim(problem.Solution, " \n")
	if id == -1 {
		problem, err := b.createProblem(problem)
		return problem.Id, err
	} else {
		err := b.updateProblem(problem)
		return problem.Id, err
	}
}

func (b *Box) ProblemSubmit(r *http.Request, user int64) (interface{}, error) {
	// Record this session
	var sess Session
	var err error
	if sess.Problem, err = strconv.ParseInt(r.FormValue("id"), 10, 64); err != nil {
		return nil, err
	}
	sess.User = user
	sess.Date = time.Now().Unix()
	sess.Code = strings.Trim(r.FormValue("code"), " ")
	if sess.Time, err = strconv.ParseInt(r.FormValue("time"), 10, 64); err != nil {
		return nil, err
	}
	sess.Solved = r.FormValue("solved") != "0"
	if err := b.storeSession(sess); err != nil {
		return nil, err
	}
	// Schedule problem for later
	n := time.Duration(b.numSuccessfulAttempts(sess.Problem, user))
	due := time.Now().Add(time.Hour*24*7*n + time.Hour)
	return nil, b.scheduleProblem(sess.Problem, user, due.Unix())
}

func (b *Box) ProblemNext(r *http.Request, user int64) (interface{}, error) {
	// Suggest the following: scheduled, not-attempted, false (write new problem)
	if p, err := b.nextScheduledProblem(user); err == nil {
		return p, err
	} else if p, err = b.notScheduledProblem(user); err == nil {
		return p, err
	} else {
		return false, nil
	}
}
