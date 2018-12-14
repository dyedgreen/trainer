package main

import (
	"log"
	"trainer/internal/pkg/server"
	"trainer/internal/pkg/auth"
)

func main() {
	s := server.New(":80")

	a := auth.New()
	defer a.Close()
	s.RegisterAuth(a)

	log.Fatal(s.ListenAndServe())
}
