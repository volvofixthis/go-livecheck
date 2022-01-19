package main

import "flag"

var configPath *string = flag.String("c", "./livechecks/livecheck.yaml", "Config file")
var metricsPath *string = flag.String("m", "", "Metrics path")
var forceStdin *bool = flag.Bool("s", false, "Force stdin input")
