package app

import (
	"github.com/hermesdt/backend-challenge/pkg/db"
)

type App struct {
	DB *db.DB
}

func New() *App {
	return &App{
		DB: db.New(),
	}
}
