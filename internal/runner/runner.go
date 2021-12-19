package runner

import (
	"fmt"
	"log"

	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
	"bitbucket.rbc.ru/go/go-livecheck/internal/validator"
	"github.com/fatih/color"
)

type Runner struct {
	Validators []validator.Validator
}

func (r *Runner) Run(data map[string]interface{}) bool {
	for _, v := range r.Validators {
		fmt.Printf("Running validator: %s ", v.Title())
		if !v.Exec(data) {
			color.Red("[Fail]\n")
			color.Blue("It's Okay to Fail, My Son\n")
			return false
		}
		color.Green("[Success]\n")
	}
	return true
}

func NewRunner(c *config.Config) *Runner {
	runner := &Runner{}
	for _, vc := range c.Validators {
		v, err := validator.NewValidator(vc)
		if err != nil {
			log.Fatalf("Problem with creating validator: %s, err: %s", vc.Title, err)
		}
		runner.Validators = append(runner.Validators, v)
	}

	return runner
}
