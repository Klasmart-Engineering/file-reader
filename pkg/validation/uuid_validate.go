package validation

import (
	"github.com/go-playground/validator/v10"
)

type ValidatedOrganizationID struct {
	Uuid string `validate:"required,uuid4"`
}
type ValidateTrackingId struct {
	Uuid string `validate:"required,uuid4"`
}

func UUIDValidate(obj interface{}) error {
	validate := validator.New()
	if c, ok := obj.(ValidatedOrganizationID); ok {
		return validate.Struct(c)
	}
	if c, ok := obj.(ValidateTrackingId); ok {
		return validate.Struct(c)
	}
	return nil
}
