package server

import (
	"encoding/json"

	uuid "github.com/satori/go.uuid"
)

type Basket struct {
	UUID       uuid.UUID   `json:"basket_id"`
	Items      []Item      `json:"items"`
	Promotions []Promotion `json:"promotions"`
}

func NewBasket() *Basket {
	uuid := uuid.NewV4()
	return &Basket{
		UUID: uuid,
	}
}

func (b *Basket) AddItem(item Item) {
	b.Items = append(b.Items, item)
}

func (b *Basket) GetTotal() int {
	total := 0
	for _, item := range b.Items {
		total += item.Price
	}

	b.Promotions = calculatePromos(b)
	for _, promo := range b.Promotions {
		total -= promo.TotalDiscount
	}

	return total
}

func (b *Basket) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		UUID       uuid.UUID   `json:"id"`
		Items      []Item      `json:"items"`
		Promotions []Promotion `json:"promotions"`
		Total      int         `json:"total"`
	}{
		b.UUID,
		b.Items,
		b.Promotions,
		b.GetTotal(),
	})
}
