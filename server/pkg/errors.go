package server

type UnkownItemError struct{}

func (e *UnkownItemError) Error() string {
	return "unkown item"
}

type BasketNotFoundError struct{}

func (e *BasketNotFoundError) Error() string {
	return "basket not found"
}

type MissingItemCodeError struct{}

func (e *MissingItemCodeError) Error() string {
	return "missin item code"
}
