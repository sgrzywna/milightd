GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
DEPCMD=dep

.PHONY: all dep build test clean

all: test dep build

dep:
	$(DEPCMD) ensure

build:
	$(MAKE) -C cmd/milightd

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	$(MAKE) clean -C cmd/milightd
