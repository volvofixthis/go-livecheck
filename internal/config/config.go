package config

import (
	"log"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

const (
	LuaEngine       = "lua"
	LuaCustomEngine = "lua_custom"
	CELEngine       = "cel"
	ScriptEngine    = "script"
	L4Engine	       	= "l4"
)

type ValidatorConfig struct {
	Type        string                 `mapstructure:"type"`
	Rule        string                 `mapstructure:"rule"`
	Title       string                 `mapstructure:"title"`
	Description string                 `mapstructure:"description"`
	Name        string                 `mapstructure:"name"`
	Major       bool                   `mapstructure:"major"` // v3 and up
	Extra       map[string]interface{} `mapstructure:"extra"`
}

type InputMetricsConfig struct {
	Type   string                 `mapstructure:"type"`
	Src    string                 `mapstructure:"src"`
	Regexp string                 `mapstructure:"regexp"` // deprecated
	Extra  map[string]interface{} `mapstructure:"extra"`
}

type OutputMetricsConfig struct {
	Type  string                 `mapstructure:"type"`
	Dst   string                 `mapstructure:"dst"`
	Extra map[string]interface{} `mapstructure:"extra"`
}

type Config struct {
	Version       string               `mapstructure:"version"`
	Validators    []*ValidatorConfig   `mapstructure:"validators"`
	InputMetrics  *InputMetricsConfig  `mapstructure:"input_metrics"`
	OutputMetrics *OutputMetricsConfig `mapstructure:"output_metrics"`
}

func GetConfig(path string) *Config {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s\n", err)
	}
	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		color.Red("Problem with unmarshalling config: %s", err)
	}
	return &config
}
