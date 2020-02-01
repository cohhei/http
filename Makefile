PHONY: test

GOFILES = $(shell ls -r *.go)

test:
	go test -v ${GOFILES}