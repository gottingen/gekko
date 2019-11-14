GOVERSION := $(shell go version | cut -d ' ' -f 3 | cut -d '.' -f 2)

.PHONY: test test-race test-cover-html help
.DEFAULT_GOAL := help

test: ## Run tests
	go test ./cast/...
	go test ./gflag/...
	go test ./peek/...


test-race: ## Run tests with race detector
	go test -race ./cast/...
	go test -race ./gflag/...
	go test -race ./peek/...

test-cover-html: ## Generate test coverage report
	go test ./cast/... -coverprofile=cast_coverage.out -covermode=count
	go test ./gflag/... -coverprofile=gflag_coverage.out -covermode=count
	go test ./peek/... -coverprofile=peek_coverage.out -covermode=count
	go tool cover -func=cast_coverage.out
	go tool cover -func=gflag_coverage.out
	go tool cover -func=peek_coverage.out

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
