.PHONY: lint test build run dev check

APP_NAME = app
CLI_NAME = artisan
BUILD_DIR = $(PWD)/build
MIGRATION_FOLDER = database/migrations/postgresql
DATABASE_URL = postgres://user:secret@localhost:5432/gfly?sslmode=disable

mod:
	go list -m --versions

all: critic security vulncheck lint test doc build

check: critic security vulncheck lint ## - Check code style, secure, lint,...

lint:
	golangci-lint run ./...

critic: ## - Check go critic
	gocritic check -enableAll -disable=unnamedResult,unlabelStmt,hugeParam,singleCaseSwitch,builtinShadow,typeAssertChain ./...

security: ## - Check go secure
	gosec -exclude-dir=core -exclude=G101,G115 ./...

vulncheck: ## - Check go vuln
	govulncheck ./...

test:
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

test.coverage:
	go tool cover -html=cover.out

build: lint test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(CLI_NAME) app/console/cli.go
	cp .env build/

run: lint test doc build
	$(BUILD_DIR)/$(APP_NAME)

start: run

schedule: build
	./build/artisan schedule:run

queue: build
	./build/artisan queue:run

migrate.up:
	migrate -path $(MIGRATION_FOLDER) -database "$(DATABASE_URL)" up

migrate.down:
	migrate -path $(MIGRATION_FOLDER) -database "$(DATABASE_URL)" down

dev:
	air -build.exclude_dir=node_modules,public,resources,Dev,bin,build,dist,docker,storage,tmp,database,docs main.go

clean:
	go mod tidy
	go clean -cache
	go clean -testcache

doc:
	swag init --parseDependency --parseInternal --parseDepth 1
	cp ./docs/swagger.json ./public/docs/

docker.run:
	docker-compose -f docker/docker-compose.yml -p gfly up -d db
	docker-compose -f docker/docker-compose.yml -p gfly up -d mail
	docker-compose -f docker/docker-compose.yml -p gfly up -d redis
	#docker-compose -f docker/docker-compose.yml -p gfly up -d minio

docker.logs:
	docker-compose -f docker/docker-compose.yml -p gfly logs -f db &
	docker-compose -f docker/docker-compose.yml -p gfly logs -f mail &
	docker-compose -f docker/docker-compose.yml -p gfly logs -f redis &
	#docker-compose -f docker/docker-compose.yml -p gfly logs -f minio &

docker.stop:
	docker-compose -f docker/docker-compose.yml -p gfly kill

docker.delete:
	docker-compose -f docker/docker-compose.yml -p gfly down

upgrade:
	go get -u all
