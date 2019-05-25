package db

import "fmt"

type RecordNotFoundError struct {
	Model string
	ID    interface{}
}

func (e *RecordNotFoundError) Error() string {
	return fmt.Sprintf("%s %s not found", e.Model, e.ID)
}
