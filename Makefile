all: deps test dist

BUILD=$(shell git rev-parse HEAD)

deps: clean
	go mod tidy
	go mod download
	go mod vendor

test: lint
	go test -race

lint: clean
	go vet ./...
	golint -set_exit_status $(shell go list ./...)

clean:
	rm -rf dist/*


dist: clean
	go build -o dist/gobunny
