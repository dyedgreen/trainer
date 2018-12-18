// The package flashcard retrieves problems
// and records a users progress

package flashcard

import (
	"math/rand"
	"time"
	"trainer/internal/pkg/problem"
)

type Card struct {
	Repeat  time.Time
	Problem problem.Problem
}

type Log struct {
	Time     time.Time
	Duration time.Duration
	DidSolve bool
	Problem  problem.Problem
}

type Box struct {
	Contains map[problem.Problem]bool
	Queue    []Card
	Log      []Log
}

func OpenBox(user string) (*Box, error) {
	// TODO
	var b Box
	b.Contains = make(map[problem.Problem]bool)
	b.Queue = make([]Card)
	b.Log = make([]Log)
	return &b, nil
}

func (b *Box) DrawNewProblem() problem.Problem {

}

func (b *Box) NextProblem() problem.Problem {

}

func (b *Box) LogProblem(prob problem.Problem) error {

}
