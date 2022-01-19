package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"regexp"

	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
	"bitbucket.rbc.ru/go/go-livecheck/internal/runner"
	"github.com/fatih/color"
)

func main() {
	flag.Parse()
	config := config.GetConfig(*configPath)
	runner, err := runner.NewRunner(config)
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
					files, err = filterFilesWithRegexp(files, config.InputMetrics.Regexp)
					if err != nil {
						os.Exit(1)
					}
					if len(files) == 0 {
						color.Red("No metrics files matching regexp: %s", config.InputMetrics.Regexp)
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
			color.Red("Can't find such input metrics type: %s", config.InputMetrics.Type)
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

func filterFilesWithRegexp(files []string, r string) ([]string, error) {
	filesF := make([]string, 0, len(files))
	for _, file := range files {
		base := filepath.Base(file)
		found, err := regexp.MatchString(r, base)
		if err != nil {
			color.Red("Regexp is wrong: %s", r)
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
