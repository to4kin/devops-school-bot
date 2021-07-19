.PHONY: build
build:
	go build -v ./cmd/devopsschoolbot

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEVAULT_GOAL := build
