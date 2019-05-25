package db

type DB struct {
	BasketsStore BasketsStore
	ItemsStore   ItemsStore
}

func New() *DB {
	return &DB{
		BasketsStore: NewBasketsMemStore(),
		ItemsStore:   NewItemsMemStore(),
	}
}
