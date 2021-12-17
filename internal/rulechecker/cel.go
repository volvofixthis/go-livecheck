package rulechecker

import (
	"errors"

	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
)

type LuaRuleChecker struct {
	config      *config.CheckerConfig
	title       string `json:"title"`
	description string `json:"description"`
}

func (rc *LuaRuleChecker) Run() bool {
	return true
}

func (rc *LuaRuleChecker) Title() string {
	return rc.config.Title
}

func NewLuaRuleChecker(config *config.CheckerConfig) (*LuaRuleChecker, error) {
	return nil, errors.New("no Lua support yet")
}
