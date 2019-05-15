package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTotalAmount(t *testing.T) {
	basket := NewBasket()
	basket.AddItem(Item{Price: 100})
	basket.AddItem(Item{Price: 1200})
	basket.AddItem(Item{Price: 50})

	assert.Equal(t, Total(13.5), basket.GetTotal())
}

func TestAddItemToBasket(t *testing.T) {
	voucher := Items[VOUCHER_CODE]
	basket := NewBasket()
	basket.AddItem(voucher)

	assert.Equal(t, 1, len(basket.Items))
	assert.Equal(t, voucher, basket.Items[0])
}
