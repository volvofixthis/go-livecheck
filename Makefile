.PHONY: build

OS=linux
ARCH=amd64

build:
	GOOS=${OS} GOARCH=${ARCH} go build -o output/go-livecheck_${OS}_${ARCH} ./cmd/livecheck

test:
	@go test ./internal/validator ./internal/runner

integration-test: build
	./livechecks/metrics.json.sh | ./output/livecheck -s -c ./livechecks/livecheck_all.yaml
	! ./livechecks/metrics.json.sh | ./output/livecheck -s -c ./livechecks/livecheck_lua.yaml
	! ./livechecks/metrics.json.sh | ./output/livecheck -s -c ./livechecks/livecheck_cel.yaml
	./livechecks/metrics.json.sh | ./output/livecheck -s -c ./livechecks/livecheck_script.yaml
	./livechecks/metrics.pid.json.sh && ./output/livecheck -c ./livechecks/livecheck_workers_cel.yaml
	./livechecks/metrics.pid.problem.json.sh && ! ./output/livecheck -c ./livechecks/livecheck_workers_cel.yaml
	./livechecks/metrics.json.sh | ./output/livecheck -s -c ./livechecks/livecheck_workers_cel.yaml
	./livechecks/metrics.json.sh | ./output/livecheck -s -c ./livechecks/livecheck_l4_tcp.yaml
	export CLIENT_CONNECTED=1 && ./output/livecheck -c ./livechecks/livecheck_cel_env.yaml -e
	./output/livecheck -c ./livechecks/livecheck_cel_v4_file.yaml
