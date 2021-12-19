package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
	"bitbucket.rbc.ru/go/go-livecheck/internal/runner"
	"github.com/spf13/viper"
)

func main() {
	configPath := flag.String("c", "./livechecks/livecheck.yaml", "Config file")
	flag.Parse()
	viper.SetConfigFile(*configPath)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %w \n", err)
	}
	config := config.Config{}
	viper.Unmarshal(&config)
	runner := runner.NewRunner(&config)
	data := map[string]interface{}{}
	err = json.NewDecoder(os.Stdin).Decode(&data)
	if err != nil {
		log.Fatalf("Error parsing metrics %s\n", err)
	}
	if !runner.Run(data) {
		os.Exit(1)
	}
}
