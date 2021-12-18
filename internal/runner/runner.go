package runner

import "bitbucket.rbc.ru/go/go-livecheck/internal/validator"

type Runner struct {
	Validators []*validator.Validator `json:"validators" yaml:"validators"`
}
