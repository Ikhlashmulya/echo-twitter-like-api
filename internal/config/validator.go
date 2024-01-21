package config

import "github.com/go-playground/validator/v10"

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(data any) error {
	if err := v.validate.Struct(data); err != nil {
		return err
	}
	return nil
}