.PHONY: build test clean

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BIN_DIR=bin
BINARY_NAME=sockmon

all: test build

build:
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME)  ./cmd/sockmon 

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
