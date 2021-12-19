package validator

import (
	"time"

	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
)

type CELValidator struct {
	config  *config.ValidatorConfig
	program cel.Program
}

func (v *CELValidator) Exec(data map[string]interface{}) (bool, error) {
	out, _, err := v.program.Eval(map[string]interface{}{
		"data": types.NewStringInterfaceMap(types.DefaultTypeAdapter, data),
		"now":  time.Now().UTC(),
	})
	if err != nil {
		return false, err
	}
	if r, ok := out.Value().(bool); ok {
		return r, nil
	}
	return false, nil
}

func (v *CELValidator) Title() string {
	return v.config.Title
}

func NewCELValidator(c *config.ValidatorConfig) (*CELValidator, error) {
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("data", decls.NewMapType(decls.String, decls.Dyn)),
			decls.NewVar("now", decls.Timestamp),
		),
	)
	if err != nil {
		return nil, err
	}
	ast, issues := env.Compile(c.Rule)
	if err := issues.Err(); issues != nil && err != nil {
		return nil, err
	}
	prg, err := env.Program(ast)
	if err != nil {
		return nil, err
	}
	return &CELValidator{config: c, program: prg}, nil
}
