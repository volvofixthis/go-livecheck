package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/volvofixthis/go-livecheck/internal/clients"
	"github.com/volvofixthis/go-livecheck/internal/config"
	"github.com/volvofixthis/go-livecheck/internal/inputmetrics"
	"github.com/volvofixthis/go-livecheck/internal/runner"

	"net/url"

	"github.com/fatih/color"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	flag.Parse()
	clients.InitHTTPClient(*insecureSkipVerify)
	config, err := config.GetConfig(*configPath, *executeTemplate, *verbose)
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
		if !runner.Run(ctx, data) {
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
					if !runner.Run(ctx, data) {
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
			if !runner.Run(ctx, data) {
				os.Exit(1)
			}
		default:
			color.Red("Can't find such input metrics type: %s", config.InputMetrics.Type)
			os.Exit(1)

		}
		return
	}
	data := map[string]interface{}{}
	if *forceStdin {
		color.Yellow("Parsing metrics from stdin")
		data, err = inputmetrics.GetMetricsStdinData(d)
		if err != nil {
			color.Red("Error reading stdin: %s", err)
			os.Exit(1)
		}
	}
	var result bool
	if *daemon {
	checkLoop:
		for {
			result = runner.Run(ctx, data)
			time.Sleep(time.Second)
			select {
			case <-ctx.Done():
				break checkLoop
			case <-time.After(*interval):
			}
		}
	} else {
		result = runner.Run(ctx, data)
	}
	if !result {
		os.Exit(1)
	}
}
