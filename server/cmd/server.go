package main

import (
	"log"
	"net/http"

	"github.com/hermesdt/backend-challenge/pkg/api"
	"github.com/hermesdt/backend-challenge/pkg/app"
)

func main() {
	app := app.New()
	api := api.New(app)

	log.Println("Starting server on port 3000")
	http.ListenAndServe(":3000", api.SetupRouter())
}
