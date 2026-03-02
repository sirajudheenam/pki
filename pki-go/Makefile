# Makefile
VERSION ?= 1.0.0
DOCKER_USER = sirajudheenam
IMAGES = server client
IMAGE_PREFIX = go-mtls-
HOSTNAME ?= localhost ## go-mtls-server-service
CERT_BASE_DIR = certs
SERVER_PORT ?= 8443
GO_FILES = $(shell find . -type f -name '*.go')
# Additional
PKI_VERSION=local
SERVER_ROOT_PATH=/hello
CERT_BASE_DIR=certs
SERVER_CERT_SUB_DIR=server
CLIENT_CERT_SUB_DIR=client

.PHONY: all build test tag push release clean cert-gen cert-clean run-server run-client help

all: test build tag

## Generate certificates for the given hostname (default: localhost)
cert-gen:
	@echo "🔐 Generating certificates for $(HOSTNAME)..."
	@mkdir -p $(CERT_BASE_DIR)/$(HOSTNAME)/{ca,server,client}
	@cd .. && bash scripts/create-certs.sh all $(HOSTNAME)
	@echo "✅ Certificates generated in $(CERT_BASE_DIR) directory"

## Remove all generated certificates
cert-clean:
	@echo "🧹 Cleaning up certificates..."
	@rm -rf $(CERT_BASE_DIR)/*
	@echo "✅ Certificates cleaned"

## Run server locally
run-server: cert-gen
	@echo "🚀 Starting server on port $(SERVER_PORT) for $(HOSTNAME)..."
	HOSTNAME=$(HOSTNAME) \
	PORT=$(SERVER_PORT) \
	CERT_BASE_DIR=$(CERT_BASE_DIR) \
	go run cmd/server/main.go

## Run client locally
run-client: cert-gen
	@echo "🔄 Starting client for $(HOSTNAME)..."
	SERVER_NAME=$(HOSTNAME) \
	SERVER_PORT=$(SERVER_PORT) \
	SERVER_ROOT_PATH=/hello \
	CLIENT_CERTS=./$(CERT_BASE_DIR)/$(HOSTNAME)/$(CLIENT_CERT_SUB_DIR) \
	go run cmd/client/main.go

## Run tests with coverage
test: cert-gen
	@echo "🧪 Running tests with coverage..."
	@mkdir -p coverage
	go test -toolexec= -v -race -coverprofile=coverage.out -covermode=atomic ./internal/... -count=1
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Tests complete. Coverage report generated in coverage.html"
	@echo "🔍 Coverage summary:"
	@go tool cover -func=coverage.out

## Run tests with coverage
test-only:
	@echo "🧪 Running tests with coverage..."
	@mkdir -p coverage
	go test -toolexec= -v -race -coverprofile=coverage.out -covermode=atomic ./internal/... -count=1
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Tests complete. Coverage report generated in coverage.html"
	@echo "🔍 Coverage summary:"
	@go tool cover -func=coverage.out

## Format code and run linters
lint:
	@echo "🔍 Running linters..."
	go fmt ./...
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi
## Run gosec
gosec:
	@echo "🔍 Running gosec..."
	@if command -v gosec >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️ gosec not installed. Run: go install github.com/securego/gosec/v2/cmd/gosec@latest"; \
	fi

## Build all images with version tag
build:
	@for img in $(IMAGES); do \
		echo "==> Building $(IMAGE_PREFIX)$$img:$(VERSION)"; \
		docker build -t $(DOCKER_USER)/go-mtls-$$img:$(VERSION) -f Dockerfile.$$img . ; \
	done

## Tag images also as :latest
tag:
	@for img in $(IMAGES); do \
		echo "==> Tagging $(IMAGE_PREFIX)$$img:latest"; \
		docker tag $(DOCKER_USER)/$(IMAGE_PREFIX)$$img:$(VERSION) $(DOCKER_USER)/$(IMAGE_PREFIX)$$img:latest ; \
	done

## Push both version and latest
push:
	@for img in $(IMAGES); do \
		echo "==> Pushing $(IMAGE_PREFIX)$$img:$(VERSION) and latest"; \
		docker push $(DOCKER_USER)/$(IMAGE_PREFIX)$$img:$(VERSION) ; \
		docker push $(DOCKER_USER)/$(IMAGE_PREFIX)$$img:latest ; \
	done

## Shortcut: build + tag + push
release: build tag push
	@echo "✅ Release complete: $(VERSION) + latest"

# ## Remove local images
clean:
	@for img in $(IMAGES); do \
		echo "==> Removing local image $$img"; \
		docker rmi -f $(DOCKER_USER)/$(IMAGE_PREFIX)$$img:$(VERSION) $(DOCKER_USER)/$(IMAGE_PREFIX)$$img:latest || true ; \
	done
	docker container prune -f
	docker image prune -f
# 	docker volume prune -f
# 	docker network prune -f

## Run services in background
up: cert-gen
	@echo "🚀 Starting services..."
	docker compose up -d

## Stop services
down:
	@echo "🛑 Stopping services..."
	docker compose down

## Tail logs
logs:
	docker compose logs -f

## Show help
help:
	@echo "Available targets:"
	@echo "  cert-gen     - Generate certificates for given HOSTNAME (default: localhost)"
	@echo "  cert-clean   - Remove all generated certificates"
	@echo "  run-server   - Run server locally (generates certs if needed)"
	@echo "  run-client   - Run client locally (generates certs if needed)"
	@echo "  test         - Run tests with coverage"
	@echo "  lint         - Format code and run linters"
	@echo "  build        - Build Docker images"
	@echo "  tag         - Tag images as latest"
	@echo "  push        - Push images to registry"
	@echo "  release     - Build, tag, and push a new release"
	@echo "  clean       - Remove images and certificates"
	@echo "  up          - Start services with Docker Compose"
	@echo "  down        - Stop services"
	@echo "  logs        - Show service logs"
	@echo
	@echo "Variables:"
	@echo "  HOSTNAME    - Target hostname for certificates (default: localhost)"
	@echo "  VERSION     - Image version tag (default: 1.0.2)"
	@echo "  SERVER_PORT - Server port (default: 8443)"
