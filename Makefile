.DEFAULT_GOAL := help

.PHONY: help test

test: ## Runs tests
	go test

run-example1: clean ## Runs sample app
	@echo "* Creating docker container with PostgreSQL"
	docker run --name umbrella-example1 -d -e POSTGRES_PASSWORD=upass -e POSTGRES_USER=uuser -e POSTGRES_DB=udb -p 54321:5432 postgres:13
	@echo "* Sleeping for 10 seconds to give database time to initialize..."
	@sleep 10
	@echo "* Building and starting application..."
	cd cmd/example1 && go build .
	cd cmd/example1 && ./example1
	@echo "* Removing previously created docker container..."
	

clean: ## Removes all created dockers
	docker rm -f umbrella-example1

help: ## Displays this help
	@awk 'BEGIN {FS = ":.*##"; printf "$(MAKEFILE_NAME)\n\nUsage:\n  make \033[1;36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[1;36m%-25s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
