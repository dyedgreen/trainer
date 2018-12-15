// The problem package can retrieve and evaluate
// coding problems

package problem

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

// Constants and erors

const (
	ProblemDir = "./data/problems"
)

var (
	ErrNotExist = errors.New("Problem not found")
)

// Problem datatype

type Problem string

type Subject string

type ProblemQuestion struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Input string `json:"input`
}

type ProblemSolution struct {
	Solution string `json:"solution"`
	Output   string `json:"output"`
	Time     string `json:"time"`
	Space    string `json:"space"`
}

func (p Problem) Subject() Subject {
	return Subject(string(p)[:strings.Index(string(p), "/")])
}

func (p Problem) Path() string {
	return ProblemDir + "/" + string(p) + ".json"
}

func (p Problem) Exists() bool {
	_, err := os.Stat(p.Path())
	return err == nil
}

func (p Problem) load(out interface{}) error {
	if !p.Exists() {
		return ErrNotExist
	}
	f, err := os.Open(p.Path())
	defer f.Close()
	if err != nil {
		return err
	}
	dec := json.NewDecoder(f)
	return dec.Decode(out)
}

func (p Problem) Question() (*ProblemQuestion, error) {
	var question ProblemQuestion
	return &question, p.load(&question)
}

func (p Problem) Solution() (*ProblemSolution, error) {
	var solution ProblemSolution
	return &solution, p.load(&solution)
}
