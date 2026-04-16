package helpers

import "github.com/go-playground/validator"

type ValidatorService struct {
	validator *validator.Validate
}

func NewValidatorService() *ValidatorService {
	return &ValidatorService{
		validator: validator.New(),
	}
}

func (v *ValidatorService) Validate(i interface{}) error {
	if v == nil || v.validator == nil {
		return nil
	}

	return v.validator.Struct(i)
}