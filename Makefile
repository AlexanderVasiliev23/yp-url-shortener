check: vet test

vet:
	go vet ./...

test:
	go test ./...

golangci-lint-run:
	docker run --rm \
        -v .:/app \
        -v ./golangci-lint/.cache/golangci-lint/v1.53.3:/root/.cache \
        -w /app \
        golangci/golangci-lint:v1.53.3 \
            golangci-lint run \
                -c .golangci.yml \
            > ./golangci-lint/report-unformatted.json
