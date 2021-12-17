package rulechecker

import (
	"errors"

	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
)

type RuleChecker interface {
	Run() bool
	Title() string
}

func NewRuleChecker(c *config.CheckerConfig) (RuleChecker, error) {
	switch c.Type {
	case config.LuaEngine:
		return NewLuaRuleChecker(c)
	}
	return nil, errors.New("Unsupported rule type")
}
