# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'



# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
.PHONY: audit
audit: test
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)" 
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage (HTML)
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## test/api: run only API (cmd/api) tests
.PHONY: test/api
test/api:
	go test -v -race -buildvcs ./cmd/api

## test/business: run only business-layer tests
.PHONY: test/business
test/business:
	go test -v -race -buildvcs ./business/...

## cover/api: API tests with coverage (func summary)
.PHONY: cover/api
cover/api:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage_api.out ./cmd/api
	go tool cover -func=/tmp/coverage_api.out

## cover/api/html: open HTML coverage for API
.PHONY: cover/api/html
cover/api/html:
	go tool cover -html=/tmp/coverage_api.out

## cover/all: all packages coverage (func summary)
.PHONY: cover/all
cover/all:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage_all.out ./...
	go tool cover -func=/tmp/coverage_all.out

## upgradeable: list direct dependencies that have upgrades available
.PHONY: upgradeable
upgradeable:
	@go run github.com/oligot/go-mod-upgrade@latest


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

## build: build the cmd/api application
.PHONY: build
build:
	go build -o=/tmp/bin/api ./cmd/api
	
## run: run the cmd/api application
.PHONY: run
run: build
	/tmp/bin/api $(ARGS)


## run/dev: start postgres container and run dev server
.PHONY: run/dev
run/dev: db/up
	sleep 2 # Give container time to initialize
	go run ./cmd/api \
		-db-dsn="postgres:postgres@localhost:5432/byfood?sslmode=disable" \
		-db-automigrate=true

		
## swagger/gen: generate Swagger docs from code annotations
.PHONY: swagger/gen
swagger/gen:
	go run github.com/swaggo/swag/cmd/swag@latest init -g ./cmd/api/main.go -o ./cmd/api/docs


# ==================================================================================== #
# SQL MIGRATIONS
# ==================================================================================== #

## migrations/new name=$1: create a new database migration
.PHONY: migrations/new
migrations/new:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest create -seq -ext=.sql -dir=./assets/migrations ${name}

## migrations/up: apply all up database migrations
.PHONY: migrations/up
migrations/up:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="postgres://${DB_DSN}" up

## migrations/down: apply all down database migrations
.PHONY: migrations/down
migrations/down:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="postgres://${DB_DSN}" down

## migrations/goto version=$1: migrate to a specific version number
.PHONY: migrations/goto
migrations/goto:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="postgres://${DB_DSN}" goto ${version}

## migrations/force version=$1: force database migration
.PHONY: migrations/force
migrations/force:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="postgres://${DB_DSN}" force ${version}

## migrations/version: print the current in-use migration version
.PHONY: migrations/version
migrations/version:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="postgres://${DB_DSN}" version


# ==================================================================================== #
# DATABASE DOCKER SETUP
# ==================================================================================== #

## db/up: start a local Postgres container
.PHONY: db/up
db/up:
	docker run -d \
		--name crud-db \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=crud \
		-p 5432:5432 \
		-v pgdata:/var/lib/postgresql/data \
		postgres:16

## db/down: stop and remove the local Postgres container
.PHONY: db/down
db/down:
	docker rm -f crud-db

## db/psql: open psql CLI inside the running container
.PHONY: db/psql
db/psql:
	docker exec -it crud-db psql -U postgres -d crud
