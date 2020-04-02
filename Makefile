NAME:=tezos

VERSION:=$(CI_COMMIT_REF_NAME)

ifeq ($(VERSION),)
	# Looks like we are not running in the CI so default to current branch or tag
	VERSION:=$(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
endif

# Need to wrap in "bash -c" so env vars work in the compiler as well as on the cli to specify the output
BUILD_CMD:=bash -c 'go build -ldflags "-X main.version=$(VERSION)" -o bin/$(NAME)-$(VERSION)-$$GOOS-$$GOARCH cmd/*'

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 $(BUILD_CMD)
	GOOS=darwin GOARCH=amd64 $(BUILD_CMD)

.PHONY: check
check: lint test

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint --enable gofmt run
