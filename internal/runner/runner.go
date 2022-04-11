package runner

import (
	"fmt"

	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
	"bitbucket.rbc.ru/go/go-livecheck/internal/outputmetrics"
	"bitbucket.rbc.ru/go/go-livecheck/internal/validator"

	"github.com/fatih/color"
)

type Runner struct {
	config        *config.Config
	Validators    []validator.Validator
	OutputMetrics outputmetrics.OutputMetrics
}

func MajorFeatureEnabled(c *config.Config) bool {
	if c.Version == "" || c.Version == "v1" || c.Version == "v2" {
		return false
	}
	return true
}

func (r *Runner) Run(data map[string]interface{}) bool {
	result := true
	for _, v := range r.Validators {
		fmt.Printf("Running validator [%s] ", v.Title())
		valid, err := v.Exec(data)

		if r.OutputMetrics != nil && v.Name() != "" {
			r.OutputMetrics.SetResult(v.Name(), valid)
		}
		if !valid {
			color.Red("[Fail]\n")
			if err != nil {
				color.Red("Validator error: %s\n", err)
			}
			if !MajorFeatureEnabled(r.config) || v.IsMajor() {
				result = false
			}
			continue
		}
		color.Green("[Success]\n")
	}
	if !result {
		color.Yellow("It's Okay to Fail, My Son\n")
	}
	if r.OutputMetrics != nil {
		r.OutputMetrics.SetResult("result", result)
		r.OutputMetrics.Flush()
	}
	return result
}

func NewRunner(c *config.Config) (*Runner, error) {
	runner := &Runner{config: c}
	for _, vc := range c.Validators {
		v, err := validator.NewValidator(vc)
		if err != nil {
			color.Red("Error creating validator [%s]: %s", vc.Title, err)
			return nil, err
		}
		runner.Validators = append(runner.Validators, v)
	}
	if c.OutputMetrics != nil {
		outputMetrics, err := outputmetrics.NewOutputMetrics(c)
		if err != nil {
			color.Red("Error setting up output metrics: %s", err)
			return nil, err
		}
		runner.OutputMetrics = outputMetrics
	}
	return runner, nil
}
