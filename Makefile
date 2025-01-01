.DEFAULT_GOAL := help

.PHONY: help test

test: ## Runs tests
	go test

run-example1: ## Runs sample app
	echo "* Creating docker container with PostgreSQL"
	docker rm -f sample-app-db
	docker run --name sample-app-db -d -e POSTGRES_PASSWORD=protopass -e POSTGRES_USER=protouser -e POSTGRES_DB=protodb -p 54320:5432 postgres:13
	echo "* Sleeping for 10 seconds to give database time to initialize..."
	sleep 10
	echo "* Building and starting application..."
	cd cmd/example1 && go build .
	cd cmd/example1 && ./example1
	echo "* Removing previously created docker container..."
	docker rm -f sample-app-db

help: ## Displays this help
	@awk 'BEGIN {FS = ":.*##"; printf "$(MAKEFILE_NAME)\n\nUsage:\n  make \033[1;36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[1;36m%-25s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
