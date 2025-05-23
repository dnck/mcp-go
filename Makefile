.PHONY: examples
help: ## The default task is help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help

examples: ## Build the examples
	mkdir -p ./bin
	go build -o ./bin/sse_server examples/sse_server/main.go
	go build -o ./bin/sse_client examples/sse_client/main.go

