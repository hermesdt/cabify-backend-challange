package main

import (
	"github.com/satori/go.uuid"
)

type Basket struct {
	UUID *uuid.UUID
	Items []Item
	State State
}

func NewBasket() *Basket {
	uuid := uuid.NewV4()
	return &Basket{
		UUID: &uuid,
	}
}