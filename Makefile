include .env
export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

init: ## Get Preliminaries
	mkdir -p bin
	go get -u github.com/valyala/quicktemplate
	go get -u github.com/valyala/quicktemplate/qtc
	go install github.com/valyala/quicktemplate/qtc@latest
	go mod tidy
	go mod download

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build: ## Build binaries, default location $(pwd)/build
	mkdir -p bin
	go build -buildvcs=false -o ./bin/mmplat ./cmd/main.go
.PHONY: build
build-win: ## Build for windows, release
	mkdir -p bin
	GOOS=windows GOARCH=amd64 go build -buildvcs=false -ldflags="-s -w -extldflags=-static" -trimpath -a -o ./bin/mmplat.exe ./cmd/main.go
qtc: ## Compile templates, default in ./internal/template
	qtc ./internal/template
