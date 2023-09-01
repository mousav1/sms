package errors

import (
	"fmt"
)

type Error struct {
	Status  int
	Message string
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error[%d] : %s", e.Status, e.Message)
}
