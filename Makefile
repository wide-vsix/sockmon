.PHONY: build test clean

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BIN_DIR=bin
BINARY_NAME=sockmon
DOCKER_IMG_TAG ?= latest

all: test build

build:
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME)  ./cmd/sockmon 

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

.PHONY: docker-build
docker-build:
	docker build -t ghcr.io/wide-vsix/sockmon:$(DOCKER_IMG_TAG) .

.PHONY: docker-push
docker-push:
	docker push ghcr.io/wide-vsix/sockmon:$(DOCKER_IMG_TAG)