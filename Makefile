dbuild:
	docker-compose build

drun:
	docker-compose up -d

ddown:
	docker-compose down -v

BINARY_NAME = awesomeProject

build:
	go build -o $(BINARY_NAME) ./cmd/logitrack

run:
	go run ./cmd/logitrack

fmt:
	go fmt ./...

lint:
	./bin/golangci-lint run ./...

test:
	go test -v ./...

test-integration:
	go test -v -tags=integration ./internal/order/...

vet:
	go vet ./...

.PHONY: start-test_db stop-test_db

start-test_db:
	docker run --name test-postgres \
		-e POSTGRES_PASSWORD=123 \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_DB=test_db \
		-p 5433:5432 \
		-d postgres

stop-test_db:
	docker stop test-postgres
	docker rm test-postgres

MIGRATE_CMD = migrate -path ./migrations -database "postgres://postgres:123@localhost:5433/test_db?sslmode=disable"

.PHONY: migrate-up migrate-down

migrate-up:
	$(MIGRATE_CMD) up

migrate-down:
	$(MIGRATE_CMD) down