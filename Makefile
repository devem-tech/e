.PHONY: test
test:
	@go test -v ./...

.PHONY: bench
bench:
	@go test -bench=. -benchmem -run=^# ./...

.PHONY: lint
lint:
	@golangci-lint run -v --fix
