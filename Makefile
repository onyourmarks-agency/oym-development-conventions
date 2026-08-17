MODULE := github.com/onyourmarks-agency/oym-development-conventions
VERSION ?= $(shell git describe --tags --always --match='v*' 2>/dev/null || echo dev)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_TIME := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

LDFLAGS := -X $(MODULE)/internal/buildinfo.Version=$(VERSION) \
	-X $(MODULE)/internal/buildinfo.GitCommit=$(GIT_COMMIT) \
	-X $(MODULE)/internal/buildinfo.BuildTime=$(BUILD_TIME)

.PHONY: build test lint clean

build:
	CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o _output/oym-conventions .

test:
	go vet ./...
	go test ./...

lint:
	golangci-lint run

clean:
	rm -rf _output
