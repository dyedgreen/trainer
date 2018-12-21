package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"trainer/internal/pkg/auth"
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

	a := auth.New(db)
	s.RegisterAuth(a)

	box := problem.NewBox(db)
	s.RegisterApiFunc("/problem/update", box.ProblemUpdate)
	s.RegisterApiFunc("/problem/submit", box.ProblemSubmit)
	s.RegisterApiFunc("/problem/next", box.ProblemNext)

	log.Fatal(s.ListenAndServe())
}
