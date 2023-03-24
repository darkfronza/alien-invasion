.DEFAULT_GOAL := build

fmt:
	go fmt ./...

lint: fmt
	golint ./...

vet: fmt
	go vet ./...

build: vet
	go build

tools: vet
	go build tools/mapgen.go

test:
	go test -v ./...

.PHONY: fmt lint vet build test