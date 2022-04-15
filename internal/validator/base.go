package validator

import (
	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
)

type ValidatorInterface interface {
	Exec(data map[string]interface{}) (bool, error)
	Title() string
	Name() string
	IsMajor() bool
}

type Validator struct {
	config *config.ValidatorConfig
}

func (v *Validator) Name() string {
	return v.config.Name
}

func (v *Validator) Title() string {
	return v.config.Title
}

func (v *Validator) IsMajor() bool {
	return v.config.Major
}

func NewValidatorBase(c *config.ValidatorConfig) *Validator {
	return &Validator{config: c}
}
