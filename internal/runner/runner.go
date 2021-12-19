package runner

import (
	"fmt"

	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
	"bitbucket.rbc.ru/go/go-livecheck/internal/validator"
	"github.com/fatih/color"
)

type Runner struct {
	Validators []validator.Validator
}

func (r *Runner) Run(data map[string]interface{}) bool {
	for _, v := range r.Validators {
		fmt.Printf("Running validator [%s] ", v.Title())
		if valid, err := v.Exec(data); !valid {

			color.Red("[Fail]\n")
			color.Yellow("It's Okay to Fail, My Son\n")
			if err != nil {
				color.Red("Validator error: %s\n")
			}
			return false
		}
		color.Green("[Success]\n")
	}
	return true
}

func NewRunner(c *config.Config) (*Runner, error) {
	runner := &Runner{}
	for _, vc := range c.Validators {
		v, err := validator.NewValidator(vc)
		if err != nil {
			color.Red("Error creating validator [%s]: %s", vc.Title, err)
			return nil, err
		}
		runner.Validators = append(runner.Validators, v)
	}

	return runner, nil
}
