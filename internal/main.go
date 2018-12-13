package main

import (
	"log"
	"trainer/internal/pkg/server"
	"trainer/internal/pkg/auth"
)

func main() {
	s := server.New(":80")

	a := &auth.Auth{}
	s.RegisterAuth(a)

	var api auth.Api
	s.RegisterApi(api)

	log.Fatal(s.ListenAndServe())
}
