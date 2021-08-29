.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

COVERAGE ?= coverage.out

.PHONY: testv0
testv0: ## Run unit tests for v0
	@echo "Running tests for v0..."
	@go test -coverprofile $(COVERAGE) ./...

.PHONY: testv2
testv2: ## Run unit tests for v2
	@echo "Running tests for v2..."
	@(cd v2 && go test -coverprofile ../$(COVERAGE) ./...)

.PHONY: test
test: ## Run unit tests for v0 and v2 combined
	@echo "mode: set" > coverage.out
	@make testv0 COVERAGE=coverage.tmp
	@tail -n +2 coverage.tmp >> coverage.out
	@make testv2 COVERAGE=coverage.tmp
	@tail -n +2 coverage.tmp >> coverage.out

.PHONY: cover
cover: ## Generate coverage report
	@make testv0
	@go tool cover -html=coverage.out
	@make testv2
	@(cd v2 && go tool cover -html=../coverage.out)
