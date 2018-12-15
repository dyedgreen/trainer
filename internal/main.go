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

	// DEBUG - test problem
	prob := problem.Problem("arrays/0")
	pa, perra := prob.Question()
	pb, perrb := prob.Solution()
	fmt.Println(prob, prob.Path(), pa, perra, pb, perrb)

	log.Fatal(s.ListenAndServe())
}
