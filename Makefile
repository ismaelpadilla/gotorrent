.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt
	go vet ./...
.PHONY:vet

golangci-lint: fmt
	golangci-lint run
.PHONY:golangci-lint

build: vet
	go build main.go
.PHONY:build
