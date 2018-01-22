.PHONY: all

GOFILES := $(shell go list -f '{{range $$index, $$element := .GoFiles}}{{$$.Dir}}/{{$$element}}{{"\n"}}{{end}}' ./... | grep -v '/vendor/')

default: clean checks test build-crossbinary

test: clean
	go test -v -cover ./...

dependencies:
	dep ensure -v

clean:
	rm -f cover.out

build:
	go build

checks: check-fmt
	gometalinter --vendor ./...

check-fmt: SHELL := /bin/bash
check-fmt:
	diff -u <(echo -n) <(gofmt -d $(GOFILES))

build-crossbinary:
	./_script/crossbinary
