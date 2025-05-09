package user

import (
	"errors"
	"fmt"
)

var ErrFirstNameRequired = errors.New("first name is required")
var ErrLastNameRequired = errors.New("last name is required")
var ErrNoFieldsToUpdate = errors.New("no fields to update")

var ErrAgeMinor = errors.New("user is not an adult")

type ErrUserNotFound struct {
	ID uint64
}

func (e ErrUserNotFound) Error() string {
	return fmt.Sprintf("user with ID %d not found", e.ID)
}
