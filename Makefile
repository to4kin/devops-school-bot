MODULE		 = $(shell env GO111MODULE=on $(GO) list -m)
DATE		?= $(shell date +%FT%T%z)
VERSION     ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || cat $(CURDIR)/.version 2> /dev/null || echo v0)
PKGS         = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
TESTPKGS     = $(shell env GO111MODULE=on $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))
BIN      	 = $(CURDIR)/bin
GO      	 = go
TIMEOUT 	 = 15

V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

.PHONY: all
all: fmt lint | $(BIN) ; $(info $(M) building executable…) @ ## Build program binary
	$Q $(GO) build \
		-tags release \
		-ldflags '-X $(MODULE)/cmd.version=$(VERSION) -X $(MODULE)/internal/app/apiserver.msgVersion=$(VERSION) -X $(MODULE)/internal/app/apiserver.msgBuildDate=$(DATE) -X $(MODULE)/cmd.buildDate=$(DATE)' \
		-o $(BIN)/$(shell basename $(MODULE)) main.go

# Tools

$(BIN):
	@mkdir -p $@

$(BIN)/%: | $(BIN) ; $(info $(M) building $(PACKAGE)…)
	$Q tmp=$$(mktemp -d); \
	   env GO111MODULE=off GOPATH=$$tmp GOBIN=$(BIN) $(GO) get $(PACKAGE) \
		|| ret=$$?; \
	   rm -rf $$tmp ; exit $$ret

GOLINT = $(BIN)/golint
$(BIN)/golint: PACKAGE=golang.org/x/lint/golint

GOCOV = $(BIN)/gocov
$(BIN)/gocov: PACKAGE=github.com/axw/gocov/...

# Tests

TEST_TARGETS := test-default test-bench test-short test-verbose test-race
.PHONY: $(TEST_TARGETS) test
test-bench:   ARGS=-run=__absolutelynothing__ -bench=. ## Run benchmarks
test-short:   ARGS=-short        ## Run only short tests
test-verbose: ARGS=-v            ## Run tests in verbose mode with coverage reporting
test-race:    ARGS=-race         ## Run tests with race detector

$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test
test: fmt lint ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests
	$Q $(GO) test -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

COVERAGE_MODE    = atomic
COVERAGE_PROFILE = $(COVERAGE_DIR)/profile.out
COVERAGE_HTML    = $(COVERAGE_DIR)/index.html
.PHONY: test-coverage test-coverage-tools
test-coverage-tools: | $(GOCOV)
test-coverage: COVERAGE_DIR := $(CURDIR)/test/coverage.$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
test-coverage: fmt lint test-coverage-tools ; $(info $(M) running coverage tests…) @ ## Run coverage tests
	$Q mkdir -p $(COVERAGE_DIR)
	$Q $(GO) test \
		-coverpkg=$$($(GO) list -f '{{ join .Deps "\n" }}' $(TESTPKGS) | \
					grep '^$(MODULE)/' | \
					tr '\n' ',' | sed 's/,$$//') \
		-covermode=$(COVERAGE_MODE) \
		-coverprofile="$(COVERAGE_PROFILE)" $(TESTPKGS)
	$Q $(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	$Q $(GO) tool cover -func=$(COVERAGE_PROFILE)

.PHONY: lint
lint: | $(GOLINT) ; $(info $(M) running golint…) @ ## Run golint
	$Q $(GOLINT) -set_exit_status $(PKGS)

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	$Q $(GO) fmt $(PKGS)


# Misc

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf $(BIN)
	@rm -rf test/coverage.*

.PHONY: help
help: ## Show this help
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'

version: ## Print version
	@echo $(VERSION)
