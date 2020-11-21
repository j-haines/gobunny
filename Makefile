all: deps test build

BUILD=$(shell git rev-parse HEAD)

build: clean
	go build -o dist/gobunny

deps: clean
	go mod tidy
	go mod download
	go mod vendor

docker: build
	docker-compose -f ./docker-compose.yml up --detach

test: lint
	go test -race $(shell go list ./...)

lint: clean
	go vet ./...
	golint -set_exit_status $(shell go list ./...)

clean:
	rm -rf dist/*

.PHONY: build deps docker test lint clean
