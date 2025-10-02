package validator

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.Validator.Struct(i)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			var errs []string
			for _, fe := range ve {
				errs = append(errs, strings.ToLower(fe.Field())+" "+fe.Tag())
			}
			return errors.New(strings.Join(errs, ", "))
		}
		return err
	}
	return nil
}
