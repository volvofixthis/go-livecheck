package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"

	"text/template"

	"bitbucket.rbc.ru/go/go-livecheck/internal/clients"
	"github.com/Masterminds/sprig/v3"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

const (
	LuaEngine       = "lua"
	LuaCustomEngine = "lua_custom"
	CELEngine       = "cel"
	ScriptEngine    = "script"
	L4Engine        = "l4"
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
	Format string                 `mapstructure:"format,omitempty"`
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

func GetConfigReader(path string) (io.Reader, error) {
	if path == "" {
		return nil, errors.New("empty config path")
	}
	switch path[0] {
	case '.', '/':
		r, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	u, err := url.ParseRequestURI(path)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case "file":
		r, err := os.Open(u.Host + u.Path)
		if err != nil {
			return nil, err
		}
		return r, nil
	case "http", "https":
		httpClient := clients.GetHTTPClient()
		resp, err := httpClient.Get(path)
		if err != nil {
			return nil, err
		}
		return resp.Body, nil
	}
	return nil, errors.New("wrong config location")
}

func GetConfig(path string, executeTemplate bool, verbose bool) (*Config, error) {
	r, err := GetConfigReader(path)

	if err != nil {
		return nil, err
	}
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	if executeTemplate {
		funcMap := sprig.TxtFuncMap()
		t, err := template.New("config").Funcs(funcMap).Parse(string(buf))
		if err != nil {
			return nil, err
		}
		bufT := &bytes.Buffer{}
		err = t.ExecuteTemplate(bufT, "config", nil)
		if err != nil {
			return nil, err
		}
		buf = bufT.Bytes()
		if verbose {
			fmt.Println("Executed template:")
			fmt.Println(string(buf))
		}
	}

	r = bytes.NewReader(buf)

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(r)
	if err != nil {
		return nil, err
	}
	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	defaults.SetDefaults(&config)
	return &config, nil
}
