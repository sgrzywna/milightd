GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

.PHONY: all build test clean

all: test build

build:
	$(MAKE) -C cmd/milightd

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	$(MAKE) clean -C cmd/milightd
