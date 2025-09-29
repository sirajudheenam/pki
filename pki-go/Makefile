# Makefile
VERSION ?= 1.0.2
DOCKER_USER = sirajudheenam
IMAGES = server client
IMAGE_PREFIX = go-mtls-

.PHONY: build test tag push release clean

# Test our code
test:
	go test -v ./internal/server -args count 1
	go test -v ./internal/client -args count 1

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
up:
	docker compose up -d

## Stop services
down:
	docker compose down

## Tail logs
logs:
	docker compose logs -f

