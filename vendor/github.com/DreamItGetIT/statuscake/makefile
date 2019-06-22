.PHONY: default deps lint test

default: deps lint test

lint:
	@golint ./...
	@go vet ./...

test:
	go test ${GOTEST_ARGS} ./...

deps:
	dep ensure