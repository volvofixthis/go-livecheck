package main

import (
	"flag"
	"time"
)

var configPath *string = flag.String("c", "./livechecks/livecheck.yaml", "Config file")
var metricsPath *string = flag.String("m", "", "Metrics path")
var forceStdin *bool = flag.Bool("s", false, "Force stdin input")
var executeTemplate *bool = flag.Bool("e", false, "Execute config as template")
var insecureSkipVerify *bool = flag.Bool("k", false, "Enable insecure skip verify")
var verbose *bool = flag.Bool("v", false, "Verbose output")
var daemon *bool = flag.Bool("d", false, "Daemon")
var interval *time.Duration = flag.Duration("i", time.Second, "Interval between checks")
