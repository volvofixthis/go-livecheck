package validator

import (
	"errors"

	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
)

type Validator interface {
	Exec(data map[string]interface{}) bool
	Title() string
}

func NewValidator(c *config.ValidatorConfig) (Validator, error) {
	switch c.Type {
	case config.LuaEngine:
		return NewLuaValidator(c)
	case config.LuaComplexEngine:
		return NewLuaValidator(c)
	case config.CELEngine:
		return NewCELValidator(c)
	}
	return nil, errors.New("unsupported rule type")
}
