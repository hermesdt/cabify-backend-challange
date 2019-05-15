package server

type Code string

type Item struct {
	Code  Code
	Name  string
	Price int
}

const (
	VOUCHER_CODE = Code("VOUCHER")
	TSHIRT_CODE  = Code("TSHIRT")
	MUG          = Code("MUG")
)

var Items = map[Code]Item{
	VOUCHER_CODE: Item{
		Code:  VOUCHER_CODE,
		Name:  "Cabify Voucher",
		Price: 500,
	},

	TSHIRT_CODE: Item{
		Code:  TSHIRT_CODE,
		Name:  "Cabify T-Shirt",
		Price: 2000,
	},

	MUG: Item{
		Code:  MUG,
		Name:  "Cabify Coffee Mug",
		Price: 750,
	},
}
