# Usecase
go-livecheck is command line tool for validation metrics and envs. If validation fails it returns exit code = 1.  
Exception for daemon mode, where go-livecheck runs forever sends metrics if they are enabled.  
For metrics validation you can use Lua or CEL. Also you can use script validator and l4 validator.
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
For testing purpose we need generate sample with metrics in json format. We will use for this simple bash script.  
Content of metrics.json.sh:  
```bash
timestamp=$(date +%s)  
timestamp=$(($timestamp-5*60))  
cat <<-EOF  
{   
    "gauge": {"client_connected": 1},  
    "timer": {"last_ping": ${timestamp}}  
}  
EOF  
```
Content of livecheck_lua.yaml:  
```yaml
validators:  
  - title: Check consumer (CEL)  
    description: Check if pool is active and worker iterated in last 10 minutes  
    type: cel  
    rule: int(data.gauge.client_connected) == 1 && (int(now) - int(data.timer.last_ping) < duration("4m").getSeconds())  
  - title: Check consumer (Lua)  
    description: Check if pool is active and worker iterated in last 10 minutes  
    type: lua  
    rule: data.gauge.client_connected == 1 and (helper:UnixTime() - data.timer.last_ping < helper:Duration("10m"))  
```
Command for validation, it will return exit code = 1 because one of validation rules will fail:  
`./livechecks/metrics.json.sh | livecheck -s -c ./livechecks/livecheck_lua.yaml`
