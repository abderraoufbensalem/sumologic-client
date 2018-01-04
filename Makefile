#!/bin/bash

.PHONY: deps vet test

deps:
	go get -u github.com/Masterminds/glide
	#glide create
	glide install

vet:
	glide nv | xargs go vet

test:
	set -o pipefail;glide nv \
		| xargs go test -v \
		| tee /dev/tty \

build:
	GOOS=linux CGO_ENABLED=0 go build -o ./main -a -ldflags '-s' -installsuffix cgo main.go