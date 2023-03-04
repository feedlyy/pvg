package helper

import (
	"errors"
)

func ValidateGet(usr string) error {
	if usr == "" {
		return errors.New("username is required")
	}

	return nil
}
