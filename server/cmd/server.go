package main

import (
	"log"

	server "github.com/hermesdt/backend-challenge/pkg"
)

func main() {
	s := server.NewServer()
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
