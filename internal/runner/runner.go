package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/volvofixthis/go-livecheck/internal/config"
	"github.com/volvofixthis/go-livecheck/internal/outputmetrics"
	"github.com/volvofixthis/go-livecheck/internal/validator"

	"github.com/fatih/color"
)

type Runner struct {
	config        *config.Config
	Validators    []validator.ValidatorInterface
	OutputMetrics outputmetrics.OutputMetrics
}

func MajorFeatureEnabled(c *config.Config) bool {
	if c.Version == "" || c.Version == "v1" || c.Version == "v2" {
		return false
	}
	return true
}

func (r *Runner) Run(ctx context.Context, data map[string]interface{}) bool {
	result := true
	vC := make(chan bool, len(r.Validators))
	for _, v := range r.Validators {
		go func(v validator.ValidatorInterface) {
			start := time.Now()

			valid, err := v.Exec(data)
			elapsed := time.Since(start)

			if r.OutputMetrics != nil && v.Name() != "" {
				r.OutputMetrics.SetResult(v.Name(), valid)
				r.OutputMetrics.SetTime(v.Name(), int64(elapsed))
			}
			fmt.Printf("Validator [%s] result ", v.Title())
			if !valid {
				color.Red("[Fail]\n")
				if err != nil {
					color.Red("Validator error: %s\n", err)
				}
				if !MajorFeatureEnabled(r.config) || v.IsMajor() {
					vC <- false
				} else {
					vC <- true
				}
			} else {
				vC <- true
				color.Green("[Success]\n")
			}
		}(v)
	}
resultCollect:
	for i := 0; i < len(r.Validators); i++ {
		select {
		case v := <-vC:
			if !v {
				result = false
			}
		case <-ctx.Done():
			result = false
			break resultCollect
		}
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
