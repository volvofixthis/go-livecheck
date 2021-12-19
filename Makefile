.PHONY: build


build:
	go build -o output/livecheck ./cmd/livecheck

test:
	@go test ./internal/validator ./internal/runner

integration-test: build
	@./livechecks/metrics.json.sh | ./output/livecheck -c ./livechecks/livecheck_all.yaml && echo "Ok!"
	@./livechecks/metrics.json.sh | ./output/livecheck -c ./livechecks/livecheck_lua.yaml || echo "Not ok!"
	@./livechecks/metrics.json.sh | ./output/livecheck -c ./livechecks/livecheck_cel.yaml || echo "Not ok!"
