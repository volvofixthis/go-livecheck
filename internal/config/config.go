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
)

type ValidatorConfig struct {
	Type        string `mapstructure:"type"`
	Rule        string `mapstructure:"rule"`
	Title       string `mapstructure:"title"`
	Description string `mapstructure:"description"`
}

type InputMetricsConfig struct {
	Type   string `mapstructure:"type"`
	Src    string `mapstructure:"src"`
	Regexp string `mapstructure:"regexp"`
}

type Config struct {
	Version      string              `mapstructure:"version"`
	Validators   []*ValidatorConfig  `mapstructure:"validators"`
	InputMetrics *InputMetricsConfig `mapstructure:"input_metrics"`
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
