package main

import (
	server "github.com/hermesdt/backend-challenge/pkg"
)

func main() {
	s := server.NewServer()
	s.Start()
}
