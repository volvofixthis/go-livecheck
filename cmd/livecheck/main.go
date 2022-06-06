package main

import (
	"flag"
	"os"
	"path/filepath"

	"bitbucket.rbc.ru/go/go-livecheck/internal/clients"
	"bitbucket.rbc.ru/go/go-livecheck/internal/config"
	"bitbucket.rbc.ru/go/go-livecheck/internal/inputmetrics"
	"bitbucket.rbc.ru/go/go-livecheck/internal/runner"

	"github.com/fatih/color"
	"net/url"
)

func main() {
	flag.Parse()
	clients.InitHTTPClient(*insecureSkipVerify)
	config, err := config.GetConfig(*configPath, *executeTemplate)
	if err != nil {
		color.Red("Error when reading config: %s", err)
		os.Exit(1)
	}
	runner, err := runner.NewRunner(config)
	if err != nil {
		color.Red("Error when creating runner: %s", err)
		os.Exit(1)
	}
	var d inputmetrics.Decoder = inputmetrics.JSONDecoder
	if config.InputMetrics != nil {
		switch config.InputMetrics.Format {
		case "yaml":
			d = inputmetrics.YAMLDecoder
		}
	}
	if *metricsPath != "" {
		data, err := inputmetrics.GetMetricsFileData(*metricsPath, d)
		if err != nil {
			os.Exit(1)
		}
		if !runner.Run(data) {
			os.Exit(1)
		}
		return
	}
	if config.InputMetrics != nil && !*forceStdin {
		src := config.InputMetrics.Src
		srcType := "file"
		if config.InputMetrics.Type != "" {
			srcType = config.InputMetrics.Type
		}
		if config.Version == "v4" {
			u, err := url.ParseRequestURI(src)
			if err != nil {
				color.Red("wrong metrics url")
				os.Exit(1)
			}
			if u.Scheme != "" {
				srcType = u.Scheme
				if srcType == "file" {
					src = u.Host + u.Path
				}
			}
		}
		switch srcType {
		case "file":
			if files, err := filepath.Glob(src); err == nil {
				if files == nil {
					color.Red("No metrics files found for pattern: %s", src)
					os.Exit(1)
				}
				if config.InputMetrics.Regexp != "" {
					files, err = inputmetrics.FilterFilesWithRegexp(files, config.InputMetrics.Regexp)
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
					data, err := inputmetrics.GetMetricsFileData(file, d)
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
		case "http", "https":
			data, err := inputmetrics.GetMetricsURLData(src, d)
			if err != nil {
				os.Exit(1)
			}
			if !runner.Run(data) {
				os.Exit(1)
			}
		default:
			color.Red("Can't find such input metrics type: %s", config.InputMetrics.Type)
			os.Exit(1)

		}
		return
	}
	data, err := inputmetrics.GetMetricsStdinData(d)
	if err != nil {
		os.Exit(1)
	}
	if !runner.Run(data) {
		os.Exit(1)
	}
}
