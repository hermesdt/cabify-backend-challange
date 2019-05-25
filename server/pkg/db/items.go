package db

import (
	"github.com/hermesdt/backend-challenge/pkg/model"
)

type ItemsStore interface {
	GetAll() map[model.Code]model.Item
}

type ItemsMemStore struct {
	Items map[model.Code]model.Item
}

func NewItemsMemStore() ItemsStore {
	return &ItemsMemStore{
		Items: make(map[model.Code]model.Item),
	}
}

func (s *ItemsMemStore) GetAll() map[model.Code]model.Item {
	return model.ITEMS
}
