# Variables
APP_NAME=auth-service
DOCKER_TAG=latest

# Go related variables.
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

# Database variables
DB_URL=postgresql://root:secret@localhost:5432/fintech_ledger?sslmode=disable

# 1. Run locally (without Docker)
run-auth:
	go run cmd/auth-service/main.go

# 2. Build the Docker image
docker-build-auth:
	docker build -f cmd/auth-service/Dockerfile -t $(APP_NAME):$(DOCKER_TAG) .

# 3. Run the Docker container
docker-run-auth:
	docker run --rm -p 8080:8080 -e SERVER_ADDRESS=:8080 $(APP_NAME):$(DOCKER_TAG)

# 4. Clean up
clean:
	rm -rf bin/
	docker rmi $(APP_NAME):$(DOCKER_TAG) || true

.PHONY: run-auth docker-build-auth docker-run-auth clean

# Database commands
db-up:
	docker-compose up -d postgres

db-down:
	docker-compose down

# Helper to verify connection
db-logs:
	docker logs -f fintech-postgres

# Database Migrations
migrate-up:
	~/go/bin/migrate -path infra/db/migration -database "$(DB_URL)" -verbose up

migrate-down:
	~/go/bin/migrate -path infra/db/migration -database "$(DB_URL)" -verbose down

# Run all
stack-up:
	docker-compose up --build -d

stack-down:
	docker-compose down

stack-logs:
	docker-compose logs -f ledger-api

test:
	go test -v ./...