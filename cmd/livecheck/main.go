package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
	"bitbucket.rbc.ru/go/go-livecheck/internal/runner"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

func main() {
	configPath := flag.String("c", "./livechecks/livecheck.yaml", "Config file")
	metricsPath := flag.String("m", "", "Metrics path")
	forceStdin := flag.Bool("s", false, "Force stdin input")
	flag.Parse()
	viper.SetConfigFile(*configPath)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s\n", err)
	}
	config := config.Config{}
	if err := viper.Unmarshal(&config); err != nil {
		color.Red("Problem with unmarshalling config: %s", err)
	}
	runner, err := runner.NewRunner(&config)
	if err != nil {
		color.Red("Error when creating runner")
		os.Exit(1)
	}
	if *metricsPath != "" {
		data, err := getMetricsFileData(*metricsPath)
		if err != nil {
			os.Exit(1)
		}
		if !runner.Run(data) {
			os.Exit(1)
		}
		return
	}
	if config.InputMetrics != nil && !*forceStdin {
		switch config.InputMetrics.Type {
		case "file":
			if files, err := filepath.Glob(config.InputMetrics.Src); err == nil {
				if files == nil {
					color.Red("No metrics files found for pattern: %s", config.InputMetrics.Src)
					os.Exit(1)
				}
				if config.InputMetrics.Regexp != "" {
					files, err = filterFIlesWithRegexp(files, config.InputMetrics.Regexp)
					if err != nil {
						os.Exit(1)
					}
					if len(files) == 0 {
						color.Red("No metrics files matched with regexp: %s", config.InputMetrics.Regexp)
						os.Exit(1)
					}
				}
				for _, file := range files {
					color.Yellow("Running validation for metrics in file: %s", file)
					data, err := getMetricsFileData(file)
					if err != nil {
						os.Exit(1)
					}
					if !runner.Run(data) {
						os.Exit(1)
					}
				}
			} else {
				color.Red("Error when searching files with metrics: %T, %s", err, err)
			}
		default:
			color.Red("No such input metrics type: %s", config.InputMetrics.Type)
			os.Exit(1)

		}
		return
	}
	data, err := getMetricsStdinData()
	if err != nil {
		os.Exit(1)
	}
	if !runner.Run(data) {
		os.Exit(1)
	}
}

func filterFIlesWithRegexp(files []string, r string) ([]string, error) {
	filesF := make([]string, 0, len(files))
	for _, file := range files {
		base := filepath.Base(file)
		found, err := regexp.MatchString(r, base)
		if err != nil {
			color.Red("Error in regexp: %s", r)
			return nil, err
		}
		if found {
			filesF = append(filesF, file)
		}
	}
	return filesF, nil
}

func getMetricsStdinData() (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := json.NewDecoder(os.Stdin).Decode(&data)
	if err != nil {
		color.Red("Can't decode json in metrics stdin: %s", err)
		return nil, err
	}
	return data, nil
}

func getMetricsFileData(path string) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	f, err := os.Open(path)
	if err != nil {
		color.Red("Can't open metrics file: %s, %s", path, err)
		return nil, err
	}
	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		color.Red("Can't decode json in metrics file: %s, %s", path, err)
		return nil, err
	}
	return data, nil
}
