package rulechecker

import (
	"errors"

	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
)

type CELRuleChecker struct {
	config      *config.CheckerConfig
	title       string `json:"title"`
	description string `json:"description"`
}

func (rc *CELRuleChecker) Run() bool {
	return true
}

func (rc *CELRuleChecker) Title() string {
	return rc.config.Title
}

func NewCELRuleChecker(config *config.CheckerConfig) (*CELRuleChecker, error) {
	return nil, errors.New("no CEL support yet")
}
