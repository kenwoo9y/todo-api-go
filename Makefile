include .env
export

.PHONY: help build-local up down logs ps migrate-mysql migrate-psql mysql psql test test-coverage lint format
.DEFAULT_GOAL := help

build-local: ## Build docker image to local development
	docker compose build --no-cache

up: ## Do docker compose up
	docker compose up

down: ## Do docker compose down
	docker compose down

logs: ## Tail docker compose logs
	docker compose logs -f

ps: ## Check container status
	docker compose ps

migrate-mysql: ## Run MySQL migration
	docker compose exec -T todo-api mysqldef -h mysql-db -P 3306 --user=$$DB_USER --password=$$DB_PASSWORD $$DB_NAME < _tools/mysql/schema.sql

migrate-psql: ## Run PostgreSQL migration
	docker compose exec todo-api psqldef -h postgresql-db -p 5432 -U $$DB_USER -W $$DB_PASSWORD $$DB_NAME -f _tools/postgresql/schema.sql

mysql: ## Access MySQL Database
	docker compose exec mysql-db mysql -u $$DB_USER -p$$DB_PASSWORD

psql: ## Access PostgreSQL Database
	docker compose exec postgresql-db psql -U $$DB_USER -d $$DB_NAME -W $$DB_PASSWORD

test: ## Run go test
	cd api && go test -v ./...

test-coverage: ## Run go test with coverage
	cd api && go test ./... -v -coverprofile=coverage.out

lint: ## Run go vet
	cd api && go vet ./...

format: ## Run go fmt
	cd api && go fmt ./...

help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'