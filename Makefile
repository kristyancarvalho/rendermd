VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: build test lint install clean

build:
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o bin/mdp ./cmd/mdp

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install ./cmd/mdp

clean:
	rm -rf bin/
