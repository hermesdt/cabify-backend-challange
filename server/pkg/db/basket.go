package db

import (
	"github.com/hermesdt/backend-challenge/pkg/model"
)

type BasketsStore interface {
	Create() model.Basket
	Get(id string) (model.Basket, error)
	Update(basket model.Basket) error
	Delete(id string)
}

type BasketsMemStore struct {
	Baskets map[string]model.Basket
}

func NewBasketsMemStore() BasketsStore {
	return &BasketsMemStore{
		Baskets: make(map[string]model.Basket),
	}
}

func (s *BasketsMemStore) Create() model.Basket {
	basket := model.NewBasket()
	s.Baskets[basket.ID.String()] = basket
	return basket
}

func (s *BasketsMemStore) Get(id string) (model.Basket, error) {
	b, ok := s.Baskets[id]
	if !ok {
		return model.Basket{}, &RecordNotFoundError{
			Model: "Basket",
			ID:    id,
		}
	}
	return b, nil
}

func (s *BasketsMemStore) Update(basket model.Basket) error {
	// TODO use a lock
	b, err := s.Get(basket.ID.String())
	if err != nil {
		return err
	}

	b.Items = make([]model.Item, 0, len(basket.Items))
	b.Items = append(b.Items, basket.Items...)

	b.Promotions = make([]model.Promotion, 0, len(basket.Promotions))
	b.Promotions = append(b.Promotions, basket.Promotions...)

	s.Baskets[b.ID.String()] = b
	return nil
}

func (s *BasketsMemStore) Delete(id string) {
	delete(s.Baskets, id)
}
