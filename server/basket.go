package main

import (
	"github.com/satori/go.uuid"
)

type Basket struct {
	UUID *uuid.UUID
	Items []Item
}

func NewBasket() *Basket {
	uuid := uuid.Must(uuid.NewV4())
	return &Basket{
		UUID: &uuid,
	}
}