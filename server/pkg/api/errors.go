package api

import "fmt"

type UnkownItemError struct{}

func (e *UnkownItemError) Error() string {
	return "unkown item"
}

type BasketNotFoundError struct {
	ID string
}

func (e *BasketNotFoundError) Error() string {
	return fmt.Sprintf("basket %s not found", e.ID)
}

type MissingItemCodeError struct{}

func (e *MissingItemCodeError) Error() string {
	return "missing item code"
}
