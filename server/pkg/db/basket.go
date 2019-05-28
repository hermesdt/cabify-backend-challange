package db

import (
	"sync"

	"github.com/hermesdt/backend-challenge/pkg/model"
)

type BasketsStore interface {
	Create() model.Basket
	Get(id string) (model.Basket, error)
	Update(basket model.Basket) error
	Delete(id string) error
}

type BasketRecord struct {
	basket *model.Basket
	locker sync.Locker
}

type BasketsMemStore struct {
	records map[string]*BasketRecord
}

func NewBasketsMemStore() BasketsStore {
	return &BasketsMemStore{
		records: make(map[string]*BasketRecord),
	}
}

func (s *BasketsMemStore) Create() model.Basket {
	basket := model.NewBasket()
	record := BasketRecord{
		basket: &basket,
		locker: &sync.Mutex{},
	}

	s.records[basket.ID.String()] = &record
	return basket
}

func (s *BasketsMemStore) Get(id string) (model.Basket, error) {
	r, err := s.get(id)
	if err != nil {
		return model.Basket{}, err
	}

	return *r.basket, nil
}

func (s *BasketsMemStore) get(id string) (*BasketRecord, error) {
	r, ok := s.records[id]
	if !ok {
		return nil, &RecordNotFoundError{
			Model: "Basket",
			ID:    id,
		}
	}
	return r, nil
}

func (s *BasketsMemStore) Update(basket model.Basket) error {
	r, err := s.get(basket.ID.String())
	if err != nil {
		return err
	}

	r.locker.Lock()
	defer r.locker.Unlock()

	b := r.basket

	b.Items = make([]model.Item, 0, len(basket.Items))
	b.Items = append(b.Items, basket.Items...)

	b.Promotions = make([]model.Promotion, 0, len(basket.Promotions))
	b.Promotions = append(b.Promotions, basket.Promotions...)

	return nil
}

func (s *BasketsMemStore) Delete(id string) error {
	r, err := s.get(id)
	if err != nil {
		return err
	}

	r.locker.Lock()
	defer r.locker.Unlock()

	delete(s.records, id)
	return nil
}
