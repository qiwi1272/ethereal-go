# Makefile for common Go commands

APP := client
PKG := ./...
BIN := bin/$(APP)

all: build

fmt:
	go fmt $(PKG)

tidy:
	go mod tidy

vet:
	go vet $(PKG)

test:
	go test -v $(PKG)

build:
	mkdir -p bin
	go build -o $(BIN)/client ./$(APP)

deps:
	go get -u ./...

all: fmt tidy vet test build

