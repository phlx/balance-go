SHELL := /bin/bash

.PHONY: all
all: run

.PHONY: help
help: ## This help overview
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

.PHONY: run
run: ## Run application in Docker Compose
	@docker-compose up -d

.PHONY: build
build: ## Build application via Docker Compose
	@docker-compose build

.PHONY: rerun
rerun: ## Stop Docker Compose, build application container, run Docker Compose
	@docker-compose down && docker-compose build app && docker-compose up -d

.PHONY: test
test: ## Run all tests inside app container
	@docker-compose exec -e CGO_ENABLED=0 app go test ./...

.PHONY: fmt
fmt: ## Run fmt with latest docker image of golang
	@docker run --rm -v $(shell pwd):/app -w /app golang:latest go fmt ./...

.PHONY: lint
lint: ## Run lint with latest image of golangci/golangci-lint
	@docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run

.PHONY: concurrent-test
concurrent-test: ## Run test for concurrent transactions in PostgreSQL, see internal/postgres/isolation.go
	@./scripts/give.sh && ./scripts/balance.sh && ./scripts/take10.sh && ./scripts/balance.sh
