package config

const (
	LuaEngine        = iota
	LuaComplexEngine = iota
	CELEngine        = iota
)

type ValidatorConfig struct {
	Type        int    `yaml:"type" json:"type"`
	Rule        string `yaml:"rule" json:"rule"`
	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description" json:"description"`
}

type Config struct {
	Validators []*ValidatorConfig `json:"validators" yaml:"validators"`
}
