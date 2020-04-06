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

.PHONY: generate
generate: install
	find sample -name "*.yml" | sed -e 's/\.yml$$//' | xargs -I% karabiner-config "%.yml" "%_generated.json"

.PHONY: emacs
emacs: generate
	karabiner-config sample/emacs.yml ~/.config/karabiner/karabiner.json

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

.PHONY: keys
keys:
	@curl -sS https://raw.githubusercontent.com/pqrs-org/Karabiner-Elements/master/src/apps/PreferencesWindow/Resources/simple_modifications.json | jq -r '.[].data.key_code' | grep -v null