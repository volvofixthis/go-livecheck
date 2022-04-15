.PHONY: build


build:
	go build -o output/livecheck ./cmd/livecheck

test:
	@go test ./internal/validator ./internal/runner

integration-test: build
	./livechecks/metrics.json.sh | ./output/livecheck -c ./livechecks/livecheck_all.yaml
	! ./livechecks/metrics.json.sh | ./output/livecheck -c ./livechecks/livecheck_lua.yaml
	! ./livechecks/metrics.json.sh | ./output/livecheck -c ./livechecks/livecheck_cel.yaml
	./livechecks/metrics.json.sh | ./output/livecheck -c ./livechecks/livecheck_script.yaml
	./livechecks/metrics.pid.json.sh && ./output/livecheck -c ./livechecks/livecheck_workers_cel.yaml
	./livechecks/metrics.pid.problem.json.sh && ! ./output/livecheck -c ./livechecks/livecheck_workers_cel.yaml
	./livechecks/metrics.pid.problem.json.sh && ./livechecks/metrics.json.sh | ./output/livecheck -s -c ./livechecks/livecheck_workers_cel.yaml
	./livechecks/metrics.json.sh | ./output/livecheck -c ./livechecks/livecheck_l4_tcp.yaml

