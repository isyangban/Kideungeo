-include .env
GOBASE=$(CURDIR)
GOBIN=$(GOBASE)/bin
GO=go

build:
	$(GO) build -o $(GOBIN) -v ./...

test:
	$(GO) test -v ./...

run-server: build
	./bin/server

.PHONY: clean
clean:
	rm -rf $(GOBIN)/*
