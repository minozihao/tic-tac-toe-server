package main

import (
	"log"
	"net/http"

	"github.com/minozihao/tic-tac-toe-server/api"
)

func main() {
	s := api.NewServer()
	log.Fatal(http.ListenAndServe(":8080", s))
}
