package main

import (
	"log"
	"net/http"
	"trainer/internal/pkg/auth"
	"trainer/internal/pkg/server"
)

type testApi int

func (t testApi) Path() string {
	return "/test"
}

func (t testApi) Call(r *http.Request, user string) (interface{}, error) {
	return user, nil
}

func main() {
	s := server.New(":80")

	a := auth.New()
	defer a.Close()
	s.RegisterAuth(a)

	test := testApi(0)
	s.RegisterApi(test)

	log.Fatal(s.ListenAndServe())
}
