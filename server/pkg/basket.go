package server

import (
	uuid "github.com/satori/go.uuid"
)

type Total float64

type Basket struct {
	UUID  *uuid.UUID
	Items []Item
}

func NewBasket() *Basket {
	uuid := uuid.NewV4()
	return &Basket{
		UUID: &uuid,
	}
}

func (b *Basket) AddItem(item Item) {
	b.Items = append(b.Items, item)
}

func (b *Basket) GetTotal() Total {
	total := 0
	for _, item := range b.Items {
		total += item.Price
	}
	return Total(float64(total) / 100)
}
