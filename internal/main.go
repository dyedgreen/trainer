package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"trainer/internal/pkg/auth"
	"trainer/internal/pkg/draft"
	"trainer/internal/pkg/problem"
	"trainer/internal/pkg/server"
)

func main() {
	s := server.New(":80")

	db, err := sql.Open("sqlite3", "./data/trainer.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	a := auth.New(db, []string{"/app/", "/api/problem/", "/api/draft/"})
	s.RegisterAuth(a)

	// Register api routes
	box := problem.NewBox(db)
	s.RegisterApiFunc("/problem/update", box.ProblemUpdate)
	s.RegisterApiFunc("/problem/submit", box.ProblemSubmit)
	s.RegisterApiFunc("/problem/next", box.ProblemNext)
	s.RegisterApiFunc("/problem/get", box.ProblemGet)

	pad := draft.NewScratchPad(db)
	s.RegisterApiFunc("/draft/update", pad.DraftUpdate)
	s.RegisterApiFunc("/draft/delete", pad.DraftDelete)
	s.RegisterApiFunc("/draft/get", pad.DraftGet)

	log.Fatal(s.ListenAndServe())
}
