package main

import (
	"fmt"
	"log"
	"net/http"
	"trainer/internal/pkg/auth"
	"trainer/internal/pkg/problem"
	"trainer/internal/pkg/server"
)

func testApi(r *http.Request, user string) (interface{}, error) {
	return user, nil
}

func main() {
	s := server.New(":80")

	a := auth.New()
	defer a.Close()
	s.RegisterAuth(a)

	s.RegisterApiFunc("/test", testApi)

	log.Fatal(s.ListenAndServe())
}
