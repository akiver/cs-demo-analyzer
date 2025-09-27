# Strip debug info
GO_FLAGS += "-ldflags=-s -w"
# Avoid embedding the build path in the executable for more reproducible builds
GO_FLAGS += -trimpath
BINARY_NAME=csda
CLI_PATH = ./cmd/cli

.DEFAULT_GOAL := help

OS = $(shell uname)
ifneq (,$(findstring MSYS_NT,$(OS)))
	IS_WINDOWS=1
endif

build-unixlike:
	@test -n "$(GOOS)" || (echo "The environment variable GOOS must be provided" && false)
	@test -n "$(GOARCH)" || (echo "The environment variable GOARCH must be provided" && false)
	@test -n "$(BIN_DIR)" || (echo "The environment variable BIN_DIR must be provided" && false)
	CGO_ENABLED=0 GOOS="$(GOOS)" GOARCH="$(GOARCH)" go build $(GO_FLAGS) -o "$(BIN_DIR)/$(BINARY_NAME)" $(CLI_PATH)
	chmod +x "$(BIN_DIR)/$(BINARY_NAME)"

build-darwin: ## Build for Darwin x64
	@"$(MAKE)" GOOS=darwin GOARCH=amd64 BIN_DIR=bin/darwin-x64 build-unixlike

build-darwin-arm64: ## Build for Darwin arm64
	@"$(MAKE)" GOOS=darwin GOARCH=arm64 BIN_DIR=bin/darwin-arm64 build-unixlike

build-linux: ## Build for Linux x64
	@"$(MAKE)" GOOS=linux GOARCH=amd64 BIN_DIR=bin/linux-x64 build-unixlike

build-linux-arm64: ## Build for Linux arm64
	@"$(MAKE)" GOOS=linux GOARCH=arm64 BIN_DIR=bin/linux-arm64 build-unixlike

build-windows: ## Build for Windows x64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(GO_FLAGS) -o bin/windows-x64/$(BINARY_NAME).exe $(CLI_PATH)

build-js: ## Build the JS bundle
	@cd js && \
	npm install && \
	npm run build

build-all: ## Run for all platform
	@"$(MAKE)" --no-print-directory -j4 \
		build-darwin \
		build-darwin-arm64 \
		build-linux \
		build-linux-arm64 \
		build-windows \
		build-js
	@cp -r ./bin/. ./js/dist/bin

npm-publish: # Publish a new version of the JS package to npm
	@test -z $(IS_WINDOWS) || (echo "Publishing from a Windows machine is not allowed because chmod would not work for unix binaries" && false)
	@test -n "$(VERSION)" || (echo "The environment variable VERSION must be provided" && false)
	@test -n "$(NPM_EMAIL)" || (echo "The environment variable NPM_EMAIL must be provided" && false)
	@npm --version > /dev/null || (echo "The npm CLI must be installed to publish" && false)
	@echo "Checking for pending git changes..." && test -z "`git status --porcelain`" || \
		(echo "Refusing to publish with these penging git changes:" && git status --porcelain && false)
	@echo "Checking for main branch..." && test "`git rev-parse --abbrev-ref HEAD`" = main || \
		(echo "Refusing to publish from non-main branch `git rev-parse --abbrev-ref HEAD`" && false)
	@echo "Checking for unpushed commits..." && git fetch
	@test "`git cherry`" = "" || (echo "Refusing to publish with unpushed commits" && false)

	@"$(MAKE)" clean
	@"$(MAKE)" build-all
	git config --global user.name AkiVer
	git config --global user.email $(NPM_EMAIL)
	@cd js && \
	npm version $(VERSION) --tag-version-prefix="" | awk '{print $$NF}' > /tmp/NEW_VERSION && \
	git add package.json package-lock.json && \
	git commit -m "chore: version `cat /tmp/NEW_VERSION`" && \
	git tag v`cat /tmp/NEW_VERSION`

	@test -z "`git status --porcelain`" || (echo "Aborting because git is somehow unclean after a commit" && false)
	@cd js && \
	npm publish --access public && \
	git push origin main --tags

publish-minor: ## Publish a minor version of the JS package
	@"$(MAKE)" VERSION=minor npm-publish

publish-patch: ## Publish a patch version of the JS package
	@"$(MAKE)" VERSION=patch npm-publish

test: ## Run all tests
	go test ./tests/ $(ARGS)

test-csgo: ## Run CS:GO tests
	go test ./tests/csgo_* $(ARGS)

test-cs2: ## Run CS2 tests
	go test ./tests/cs2_* $(ARGS)

test-verbose: ## Run tests in verbose
	@"$(MAKE)" --no-print-directory ARGS=-v test

vet: ## Run go vet
	go vet ./cmd/... ./pkg/...

clean: ## Clean up project files
	rm -rf bin
	rm -rf ./js/dist
	rm -f ./cs-demos/*.csv
	rm -f ./cs-demos/*.json
	rm -f ./*.csv

help:
	@echo 'Targets:'
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {printf "\033[36m  %-20s\033[0m %s\n", $$1, $$NF}' $(MAKEFILE_LIST)
