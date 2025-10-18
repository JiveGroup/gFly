.PHONY: lint test build run dev check help stop

help: ## Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

# Build
APP_NAME = app
CLI_NAME = artisan
BUILD_DIR = $(PWD)/build

# Database migration (Check file .env)
DB_NAME=gfly
DB_USERNAME=user
DB_PASSWORD=secret
DB_SSL_MODE=disable
MIGRATION_FOLDER = database/migrations/postgresql
DATABASE_URL = postgres://${DB_USERNAME}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=${DB_SSL_MODE}

all: check test doc build ## - Check code style, secure, lint, test, doc and build

check: critic security vulncheck lint ## - Check code style, secure, lint,...

lint: ## - Run golangci-lint to check code quality
	golangci-lint run ./...

critic: ## - Check go critic
	gocritic check -enableAll -disable=unnamedResult,unlabelStmt,hugeParam,singleCaseSwitch,builtinShadow,typeAssertChain ./...

security: ## - Check go secure
	gosec -exclude-dir=core -exclude=G101,G115 ./...

vulncheck: ## - Check go vuln
	govulncheck ./...

test: ## - Run tests with coverage report
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

test.coverage: ## - Open HTML coverage report in browser
	go tool cover -html=cover.out

build: lint test ## - Build the application and CLI tool
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) cmd/web/main.go
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(CLI_NAME) cmd/console/main.go
	cp .env build/

run: lint test doc build ## - Run the application after building
	$(BUILD_DIR)/$(APP_NAME)

start: run ## - Alias for run command

stop: ## - Stop server process on port 7889
	bash scripts/stop.sh

schedule: build ## - Run scheduled tasks
	./build/artisan schedule:run

queue: build ## - Run queue workers
	./build/artisan queue:run

migrate.up: ## - Run database migrations up
	migrate -path $(MIGRATION_FOLDER) -database "$(DATABASE_URL)" up

migrate.down: ## - Revert the last database migration
	migrate -path $(MIGRATION_FOLDER) -database "$(DATABASE_URL)" down

dev: ## - Run application in development mode with hot reloading
	air -build.exclude_dir=node_modules,public,resources,Dev,bin,build,dist,docker,storage,tmp,database,docs cmd/web/main.go

clean: ## - Clean up Go modules, cache, and test cache
	go mod tidy
	go clean -cache
	go clean -testcache

doc: ## - Generate API documentation using Swag
	swag init --parseDependency --parseDepth 1 --exclude build,database,deployments,docs,node_modules,public,resources,storage,tmp,vendor -g cmd/web/main.go
	cp ./docs/swagger.json ./public/docs/

container.run: ## - Start required Docker containers (PostgreSQL, Mail, Redis)
	docker compose --env-file deployments/container.env -f deployments/docker/docker-compose.yml -p gfly up -d db
	docker compose --env-file deployments/container.env -f deployments/docker/docker-compose.yml -p gfly up -d mail
	docker compose --env-file deployments/container.env -f deployments/docker/docker-compose.yml -p gfly up -d redis
	#docker compose --env-file deployments/container.env -f deployments/docker/docker-compose.yml -p gfly up -d minio

container.logs: ## - Show logs from all Docker containers
	docker compose --env-file deployments/container.env -f deployments/docker/docker-compose.yml -p gfly logs -f db &
	docker compose --env-file deployments/container.env -f deployments/docker/docker-compose.yml -p gfly logs -f mail &
	docker compose --env-file deployments/container.env -f deployments/docker/docker-compose.yml -p gfly logs -f redis &
	#docker compose --env-file deployments/container.env -f deployments/docker/docker-compose.yml -p gfly logs -f minio &

container.stop: ## - Stop all Docker containers
	docker compose --env-file deployments/container.env -f deployments/docker/docker-compose.yml -p gfly kill

container.delete: ## - Stop and remove all Docker containers
	docker compose --env-file deployments/container.env -f deployments/docker/docker-compose.yml -p gfly down

upgrade: ## - Upgrade all Go dependencies to their latest versions
	go get -u all

api.scripts: ## - Generate API shell scripts from Swagger file
	./scripts/generate_api_scripts.sh
