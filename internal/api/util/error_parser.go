package util

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ErrorParser(err error) error {
	// To show the proper message to the caller walk validator path
	var ve validator.ValidationErrors
	var ue *json.UnmarshalTypeError

	switch {
	case errors.As(err, &ve):
		for _, fe := range ve {
			switch fe.Tag() {
			case "required":
				return fmt.Errorf("%s field is required", fe.Field())
			}
		}
	case errors.As(err, &ue):
		return fmt.Errorf("%s field has invalid data type", ue.Field)
	}

	return err
}
