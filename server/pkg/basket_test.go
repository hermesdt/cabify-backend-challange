package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasketVoucherPromition(t *testing.T) {
	b := NewBasket()
	b.Items = []Item{
		ITEMS[VOUCHER_CODE],
		ITEMS[VOUCHER_CODE],
		ITEMS[TSHIRT_CODE],
		ITEMS[VOUCHER_CODE],
		ITEMS[MUG_CODE],
		ITEMS[VOUCHER_CODE],
		ITEMS[VOUCHER_CODE],
	}

	assert.Len(t, b.Promotions, 0)
	expectedValue := ITEMS[VOUCHER_CODE].Price*3 + ITEMS[TSHIRT_CODE].Price + ITEMS[MUG_CODE].Price
	assert.Equal(t, expectedValue, b.GetTotal())
	assert.Len(t, b.Promotions, 1)
	assert.Equal(t, b.Promotions[0].TotalDiscount, 1000)
}

func TestBasketVoucherAndTshirtPromitions(t *testing.T) {
	b := NewBasket()
	b.Items = []Item{
		ITEMS[MUG_CODE],
		ITEMS[VOUCHER_CODE],
		ITEMS[VOUCHER_CODE],
		ITEMS[TSHIRT_CODE],
		ITEMS[TSHIRT_CODE],
		ITEMS[MUG_CODE],
		ITEMS[TSHIRT_CODE],
		ITEMS[VOUCHER_CODE],
		ITEMS[TSHIRT_CODE],
		ITEMS[MUG_CODE],
	}

	assert.Len(t, b.Promotions, 0)
	expectedValue := 0 +
		ITEMS[MUG_CODE].Price*3 +
		ITEMS[VOUCHER_CODE].Price*2 +
		ITEMS[TSHIRT_CODE].Price*4 - 4*100
	assert.Equal(t, expectedValue, b.GetTotal())
}
