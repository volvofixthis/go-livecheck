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
