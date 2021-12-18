package validator

import (
	"log"
	"time"

	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

type CELValidator struct {
	config  *config.ValidatorConfig
	program cel.Program
}

func (v *CELValidator) Exec(data map[string]interface{}) bool {
	out, _, err := v.program.Eval(map[string]interface{}{
		"data": data,
		"now":  time.Now().UTC(),
	})
	if err != nil {
		log.Fatalf("eval error %s", err)
	}
	if r, ok := out.Value().(bool); ok {
		return r
	}
	return false
}

func (v *CELValidator) Title() string {
	return v.config.Title
}

func NewCELValidator(c *config.ValidatorConfig) (*CELValidator, error) {
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("data", decls.Any),
			decls.NewVar("now", decls.Timestamp),
		),
	)
	if err != nil {
		log.Fatalf("can'te create CEL env %s", err)
	}
	ast, issues := env.Compile(c.Rule)
	if issues != nil && issues.Err() != nil {
		log.Fatalf("type-check error: %s", issues.Err())
	}
	prg, err := env.Program(ast)
	if err != nil {
		log.Fatalf("program construction error: %s", err)
	}
	return &CELValidator{config: c, program: prg}, nil
}
