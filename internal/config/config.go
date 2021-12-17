package config

const (
	LuaEngine = iota
	CELEngine = iota
)

type CheckerConfig struct {
	Type        int    `yaml:"type" json:"type"`
	Rule        string `yaml:"rule" json:"rule"`
	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description" json:"description"`
}

type Config struct {
	Checkers []*CheckerConfig `yaml:"checkers" json:"checkers"`
}
