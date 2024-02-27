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

mem-optimization-diff:
	go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof