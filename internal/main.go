package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"trainer/internal/pkg/auth"
	"trainer/internal/pkg/server"
)

func testApi(r *http.Request, user string) (interface{}, error) {
	return user, nil
}

func main() {
	s := server.New(":80")

	db, err := sql.Open("sqlite3", "./data/trainer.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	a := auth.New(db)
	s.RegisterAuth(a)

	s.RegisterApiFunc("/test", testApi)

	log.Fatal(s.ListenAndServe())
}
