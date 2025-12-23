package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateReqAuthBody[T any](body T) []string {
	if err := validate.Struct(body); err != nil {
		errs := err.(validator.ValidationErrors)
		var msg []string
		for _, e := range errs {
			msg = append(msg, fmt.Sprintf(" campo %s: fallo en-> %s", e.Field(), e.Tag()))
		}
		return msg
	}
	return nil
}
