.PHONY: help checkfmt fmt vet dl-modules release pkg test clean

.DEFAULT_GOAL := help

# Root directory of the project (absolute path). Has a trailing slash.
ROOT_DIR=$(dir $(abspath $(lastword $(MAKEFILE_LIST))))
# Output directories
BIN_DIR := $(ROOT_DIR)bin
PKG_DIR := $(ROOT_DIR)pkg
# Get version information
#GIT_COMMIT := $(shell git rev-parse --short=10 HEAD)
#DESCRIBE := $(subst -, ,$(subst release/,,$(shell git describe --match 'release/*' --abbrev=10 --dirty=+m)))
# Build date
BUILT_ON := $(shell date -u '+%Y-%m-%d')
# Set Go vars
GO_SRC = $(shell go list ./... | grep -v vendor)
GOFMT_FILES = $$(find . -name '*.go' | grep -v vendor)
GO_TEST_FLAGS = -v -covermode=count -coverprofile=coverage.out
CGO_ENABLED = 0
# List all files that haven't passed gofmt
UNFMT_FILES := $(shell gofmt -l $(GOFMT_FILES))
# Do not update go.mod when building -- if the file is out of date or needs
# updating, an error will be generated.
GOFLAGS = "-mod=readonly"
GOOS = linux
GOARCH = amd64

# For use by humans
all: fmt vet test

checkfmt: ## Verifies that all files pass `go fmt`
		@echo "+ $@"
ifdef UNFMT_FILES
		@echo "Unformatted files found; please run 'make fmt' to fix:"
		@echo "$(UNFMT_FILES)"
		@exit 1
endif

fmt: ## Verifies that all files pass `go fmt`
		@echo "+ $@"
		gofmt -l -w $(GOFMT_FILES)

vet: ## Runs go vet on all non-vendor source code
		@echo "+ $@"
		go vet $(GO_SRC)

lint: ## Runs golangci-lint
		@echo "+ $@"
		golangci-lint run -v

dl-modules: ## Downloads Go modules
		@echo "+ $@"
		@echo "==> Downloading Go modules..."
		go mod download

release: dl-modules ## Builds the executable
		@echo "+ $@"
		go build -o $(BIN_DIR)/$(GOOS)/$(GOARCH)/$(NAME) $(GO_LDFLAGS) .

pkg: release ## Creates the tarball package
		@echo "+ $@"
		mkdir -p $(PKG_DIR); \
		cd $(BIN_DIR)/$(GOOS)/$(GOARCH); \
		sha256sum "$(NAME)" > "$(NAME)".sha256; \
		tar --exclude=*.sha256 -czf "$(PKG_DIR)/$(NAME)_$(VERSION)_$(GOOS)_$(GOARCH).tar.gz" ./*; \
		cd $(ROOT_DIR)

test: dl-modules ## Runs the Go unit tests with -short
		@echo "+ $@"
		go test -short $(GO_TEST_FLAGS) $(GO_SRC)

clean: ## Removes all build and test output
		@echo "+ $@"
		rm -rf $(BIN_DIR) $(PKG_DIR)
		rm -f coverage* unit_tests*
		go clean -testcache $(GO_SRC)

help: ## Print this message and exit.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%20s\033[0m : %s\n", $$1, $$2}' $(MAKEFILE_LIST)

