package config

const (
	LuaEngine       = "lua"
	LuaCustomEngine = "lua_custom"
	CELEngine       = "cel"
)

type ValidatorConfig struct {
	Type        string `yaml:"type" json:"type"`
	Rule        string `yaml:"rule" json:"rule"`
	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description" json:"description"`
}

type Config struct {
	Validators []*ValidatorConfig `json:"validators" yaml:"validators"`
}
