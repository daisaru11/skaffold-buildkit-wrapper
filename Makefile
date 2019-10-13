export GO111MODULE=on

build:
	go build

install:
	go install

test:
	go test ./... -v


.PHONY: build install test