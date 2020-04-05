BIN := build
CURRENT_REVISION := $(shell git rev-parse --short HEAD)
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null)
MODULE := "github.com/uphy/karabiner-config"
BUILD_LDFLAGS := "-X $(MODULE)/app.revision=$(CURRENT_REVISION) -X $(MODULE)/app.version=$(VERSION) -X main.a=bbb"

.PHONY: all
all: clean build

.PHONY: install
install:
	GOOS=darwin GOARCH=amd64 go install -ldflags=$(BUILD_LDFLAGS)

.PHONY: build
build:
	GOOS=darwin GOARCH=amd64 go build -ldflags=$(BUILD_LDFLAGS) -o $(BIN)/karabiner-config; \
	cd build && gzip -f karabiner-config

.PHONY: clean
clean: 
	rm -rf $(BIN)
	go clean

.PHONY: test
test:
	go test -v ./...