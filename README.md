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
