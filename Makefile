check: vet test

vet:
	go vet ./...

test-with-coverage: test coverage

test:
	go test -coverprofile=coverage.out ./...

coverage:
	 go tool cover -html=coverage.out -o coverage.html

golangci-lint-run:
	docker run --rm \
        -v .:/app \
        -v ./golangci-lint/.cache/golangci-lint/v1.53.3:/root/.cache \
        -w /app \
        golangci/golangci-lint:v1.53.3 \
            golangci-lint run \
                -c .golangci.yml \
            > ./golangci-lint/report-unformatted.json

save-base-profile:
	curl http://127.0.0.1:8081/debug/pprof/heap > ./profiles/base.pprof

save-result-profile:
	curl http://127.0.0.1:8081/debug/pprof/heap > ./profiles/result.pprof

show-base-profile:
	go tool pprof -http=":9091" -seconds=30 ./profiles/base.pprof

show-result-profile:
	go tool pprof -http=":9092" -seconds=30 ./profiles/result.pprof

base-result-profile-diff:
	go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof

format-code:
	goimports -local "github.com/AlexanderVasiliev23/yp-url-shortener" -w ./..

multichecker:
	go run cmd/staticlint/main.go ./...

run-shortener-with-build-flags:
	go run -ldflags "-X main.buildVersion=1.2.3 -X 'main.buildDate=$(date +'%Y-%m-%d %H:%M:%S')' -X 'main.buildCommit=$(git rev-parse --short HEAD)'" ./cmd/shortener/.