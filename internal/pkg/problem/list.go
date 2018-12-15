// List subjects and problems

package problem

import (
	"os"
	"strings"
)

func ListSubjects() ([]Subject, error) {
	f, err := os.Open(ProblemDir)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	entries, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	subjects := make([]Subject, 0)
	for _, e := range entries {
		if e.IsDir() {
			subjects = append(subjects, Subject(e.Name()))
		}
	}
	return subjects, nil
}

func (s Subject) ListProblems() ([]Problem, error) {
	f, err := os.Open(ProblemDir + "/" + string(s))
	defer f.Close()
	if err != nil {
		return nil, err
	}
	entries, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	problems := make([]Problem, 0)
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".json") {
			problems = append(problems, Problem(string(s) + "/" + strings.TrimSuffix(e.Name(), ".json")))
		}
	}
	return problems, nil
}
