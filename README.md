# Usecase
go-livecheck is command line tool for validation metrics and envs. If validation fails it returns exit code = 1.  
Exception for daemon mode, where go-livecheck runs forever sends metrics if they are enabled
# Install
Run go install:  
go install github.com/volvofixthis/go-livecheck/cmd/livecheck@latest 
# Build
Use make:  
make OS=linux GOARCH=arm64  
make OS=linux GOARCH=amd64  
make OS=darwin GOARCH=arm64  
make OS=darwin GOARCH=amd64  
make OS=windows GOARCH=amd64  
Or just build all platforms in one shot:  
make all
# Metrics validation
Metrics data can be passed via stdin. Data format can be yaml or json.
For test data we need generate json data, we use for this simple bash script.  
Content of metrics.json.sh:  
`
timestamp=$(date +%s)
timestamp=$(($timestamp-5*60))
cat <<-EOF
{
    "gauge": {"client_connected": 1},
    "timer": {"last_ping": ${timestamp}}
}
EOF
`
Command for validation, it will return exit code = 1 because one of validation rules fail:  
`./livechecks/metrics.json.sh | ./output/${TEST_BINARY} -s -c ./livechecks/livecheck_lua.yaml`
