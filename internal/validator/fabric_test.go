package validator

import (
	"testing"
	"time"

	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
)

func TestLuaValidatorFabric(t *testing.T) {
	c := config.ValidatorConfig{
		Type:        config.LuaEngine,
		Title:       "Check consumer (Lua)",
		Description: "Check if pool is active and worker iterated in last 10 minutes",
		Rule:        `data.client_connected and (helper:UnixTime() - data.last_livecheck < helper:Duration("10m"))`,
	}
	data := map[string]interface{}{
		"client_connected": true,
		"last_livecheck":   float64(time.Now().UTC().UnixNano())/float64(time.Second) - 5*60,
	}
	luaValidator, _ := NewLuaValidator(&c)
	if valid, err := luaValidator.Exec(data); !valid || err != nil {
		t.Error("Not valid")
	}
}

func TestCELValidatorFabric(t *testing.T) {
	c := config.ValidatorConfig{
		Type:        config.CELEngine,
		Title:       "Check consumer (CEL)",
		Description: "Check if pool is active and worker iterated in last 10 minutes",
		Rule:        `data.client_connected == true && (int(now) - int(data.last_livecheck) < duration("10m").getSeconds())`,
	}
	data := map[string]interface{}{
		"client_connected": true,
		"last_livecheck":   float64(time.Now().UTC().UnixNano())/float64(time.Second) - 5*60,
	}
	celValidator, _ := NewCELValidator(&c)
	if valid, err := celValidator.Exec(data); !valid || err != nil {
		t.Error("Not valid")
	}
}
