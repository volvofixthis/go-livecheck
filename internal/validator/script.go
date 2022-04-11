package validator

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
)

type ScriptValidator struct {
	config *config.ValidatorConfig
}

func (v *ScriptValidator) Exec(data map[string]interface{}) (bool, error) {
	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0755)
	if err != nil {
		return false, err
	}
	buf, err := json.Marshal(data)
	if err != nil {
		return false, nil
	}
	stdin := strings.NewReader(string(buf))
	ruleParts := strings.Split(v.config.Rule, " ")
	cmd := &exec.Cmd{
		Path:   ruleParts[0],
		Args:   ruleParts,
		Stdin:  stdin,
		Stdout: devnull,
		Stderr: devnull,
	}

	err = cmd.Start()
	if err != nil {
		return false, err
	}

	err = cmd.Wait()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (v *ScriptValidator) Name() string {
	return v.config.Name
}

func (v *ScriptValidator) Title() string {
	return v.config.Title
}

func (v *ScriptValidator) IsMajor() bool {
	return v.config.Major
}

func NewScriptValidator(c *config.ValidatorConfig) (*ScriptValidator, error) {
	return &ScriptValidator{config: c}, nil
}
