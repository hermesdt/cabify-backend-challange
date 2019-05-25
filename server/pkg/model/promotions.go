package model

type Promotion struct {
	Name          string `json:"name"`
	TotalDiscount int    `json:"total_discount"`
}

func calculatePromos(b *Basket) []Promotion {
	var promotions []Promotion
	if promo := buy2pay1VoucherPromo(b); promo != nil {
		promotions = append(promotions, *promo)
	}
	if promo := more3TshirtPromo(b); promo != nil {
		promotions = append(promotions, *promo)
	}

	return promotions
}

func buy2pay1VoucherPromo(b *Basket) *Promotion {
	vouchersCount := 0
	for _, item := range b.Items {
		if item.Code == VOUCHER_CODE {
			vouchersCount++
		}
	}
	if vouchersCount < 2 {
		return nil
	}

	return &Promotion{
		Name:          "Buy 2 vouchers get 1 free",
		TotalDiscount: ITEMS[VOUCHER_CODE].Price * (vouchersCount / 2),
	}
}

func more3TshirtPromo(b *Basket) *Promotion {
	tshirtCount := 0
	for _, item := range b.Items {
		if item.Code == TSHIRT_CODE {
			tshirtCount++
		}
	}
	if tshirtCount < 3 {
		return nil
	}

	return &Promotion{
		Name:          "Buy 3 or more T-Shirt for 19€ each",
		TotalDiscount: (ITEMS[TSHIRT_CODE].Price - 1900) * tshirtCount,
	}
}
