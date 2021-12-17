package rulegroup

import "bitbucket.rbc.ru/go/go-livecheck/internal/rulechecker"

type RuleGroup struct {
	Checkers []*rulechecker.RuleChecker `json:"checkers"`
}
