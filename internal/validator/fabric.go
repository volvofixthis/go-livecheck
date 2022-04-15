package validator

import (
	"errors"

	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
)

func NewValidator(c *config.ValidatorConfig) (ValidatorInterface, error) {
	switch c.Type {
	case config.LuaEngine:
		return NewLuaValidator(c)
	case config.LuaCustomEngine:
		return NewLuaValidator(c)
	case config.CELEngine:
		return NewCELValidator(c)
	case config.ScriptEngine:
		return NewScriptValidator(c)
	case config.L4Engine:
		return NewL4Validator(c)
	}
	return nil, errors.New("unsupported rule type")
}
