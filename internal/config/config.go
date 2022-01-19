package config

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
