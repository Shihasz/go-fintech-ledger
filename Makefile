# Variables
APP_NAME=auth-service
DOCKER_TAG=latest

# Go related variables.
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

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