MODULE		 = $(shell env GO111MODULE=on $(GO) list -m)
DATE		?= $(shell date +%FT%T%z)
VERSION     ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || cat $(CURDIR)/.version 2> /dev/null || echo v0)
PKGS         = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
TESTPKGS     = $(shell env GO111MODULE=on $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))
BIN      	 = $(CURDIR)/bin
GO      	 = go
TIMEOUT 	 = 15
COVERAGE_DIR = coverage

V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

.PHONY: all
all: fmt lint | $(BIN) ; $(info $(M) building executable…) @ ## build program binary
	$Q $(GO) build \
		-tags release \
		-ldflags '-X $(MODULE)/cmd.Version=$(VERSION) -X $(MODULE)/cmd.BuildDate=$(DATE)' \
		-o $(BIN)/$(basename $(MODULE)) ./cmd/devopsschoolbot/main.go

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

default: help

test: ## run test with race
	go test -v -race -timeout 30s ./...

coverage: ## Run test coverage and generate html report
	rm -fr coverage
	mkdir coverage
	go test -covermode count -coverprofile $(COVERAGE_DIR)/coverage.out ./...
	go tool cover -func=$(COVERAGE_DIR)/coverage.out
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html



.PHONY: lint
lint: | $(GOLINT) ; $(info $(M) running golint…) @ ## Run golint
	$Q $(GOLINT) -set_exit_status $(PKGS)

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	$Q $(GO) fmt $(PKGS)


# Misc

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## cleanup everything
	@rm -rf $(BIN)
	@rm -rf $(COVERAGE_DIR)

.PHONY: help
help: ## show this help
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'

version: ## print version
	@echo $(VERSION)