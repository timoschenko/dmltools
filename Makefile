MKDIRS?=mkdir -p

GO?=go
JAVA?=java
CURL?=curl
PYTHON=python3
PIP=./venv/bin/pip3
WS_FIXER=./venv/bin/trailing-whitespace-fixer
EOF_FIXER=./venv/bin/end-of-file-fixer
GODOC?=$(shell command -v godoc 2>/dev/null || echo $(shell $(GO) env GOPATH)/bin/godoc)

BUILD_DIR:=build
BUILD_DIR_FULLPATH:=$(realpath .)/$(BUILD_DIR)

BUILD_ENV_DIR:=$(BUILD_DIR)/env

# Go Settings {{{

DEFAULT_GOCACHE := $(BUILD_DIR_FULLPATH)/cache/golang
# Check if `GOCACHE` env variable exists, if no use `DEFAULT_GOCACHE`
GOCACHE?=$(DEFAULT_GOCACHE)
export GOCACHE

CI          ?= false
SKIP_ANTLR  ?= false
GIT_VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE  ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GO_LDFLAGS  = -X main.Version=$(GIT_VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildDate=$(BUILD_DATE)

# }}}

ANTLR_VERSION := 4.13.2
ANTLR_JAR = $(BUILD_DEPS_DIR)/antlr-$(ANTLR_VERSION)-complete.jar
ANTLR4_FORMATTER = 1.2.1
ANTLR_FORMAT_JAR = $(BUILD_DEPS_DIR)/antlr4-formatter-standalone-$(ANTLR4_FORMATTER).jar

G4_DML_DIR := ./grammar/dml
G4_DML := $(G4_DML_DIR)/dml.g4
BUILD_DEPS_DIR := $(BUILD_DIR)/deps

$(BUILD_DEPS_DIR):
	$(MKDIRS) $(BUILD_DEPS_DIR)

$(BUILD_ENV_DIR):
	$(MKDIRS) $(BUILD_ENV_DIR)

# Pip setup {{{
$(PIP):
	$(PYTHON) -m venv ./venv

./venv/bin/pre-commit $(WS_FIXER): $(PIP)
	$(PIP) install pre-commit pre-commit-hooks

# }}}


# Download the ANTLR jar file
$(ANTLR_JAR): | $(BUILD_DEPS_DIR)
	$(CURL) -o $@ https://www.antlr.org/download/antlr-$(ANTLR_VERSION)-complete.jar

$(ANTLR_FORMAT_JAR): | $(BUILD_DEPS_DIR)
	$(CURL) -o $@ https://repo1.maven.org/maven2/com/khubla/antlr4formatter/antlr4-formatter-standalone/$(ANTLR4_FORMATTER)/antlr4-formatter-standalone-$(ANTLR4_FORMATTER).jar

BUILD_APPS:=
BUILD_APPS+=dmlast
BUILD_APPS+=dmlfmt

# Intermediate ANTLR generated files in grammar/dml
G4_GENERATED_FILES:=
G4_GENERATED_FILES+=$(G4_DML_DIR)/dml_base_listener.go
G4_GENERATED_FILES+=$(G4_DML_DIR)/dml_lexer.go
G4_GENERATED_FILES+=$(G4_DML_DIR)/dml_listener.go
G4_GENERATED_FILES+=$(G4_DML_DIR)/dml_parser.go

# Copied files in pkg/dmlex
DMLEX_GENERATED_FILES:=
DMLEX_GENERATED_FILES+=./pkg/dmlex/dml_base_listener.go
DMLEX_GENERATED_FILES+=./pkg/dmlex/dml_lexer.go
DMLEX_GENERATED_FILES+=./pkg/dmlex/dml_listener.go
DMLEX_GENERATED_FILES+=./pkg/dmlex/dml_parser_dml_parser.go

G4_FORMAT_FLAG:=$(BUILD_ENV_DIR)/g4-format.flag
$(G4_FORMAT_FLAG): $(BUILD_ENV_DIR) $(ANTLR_FORMAT_JAR) $(G4_DML) $(WS_FIXER)
$(G4_FORMAT_FLAG):
	$(JAVA) -jar $(ANTLR_FORMAT_JAR) --input=$(G4_DML)
	-@$(WS_FIXER) $(G4_DML) || true
	-@$(EOF_FIXER) $(G4_DML) || true
	$(file > $@,)

format: $(G4_FORMAT_FLAG)

G4_GENERATED_FLAG:=$(BUILD_ENV_DIR)/g4-generated.flag
$(G4_GENERATED_FLAG): $(G4_FORMAT_FLAG) $(ANTLR_JAR) $(BUILD_ENV_DIR)
$(G4_GENERATED_FLAG):
	cd $(G4_DML_DIR) && $(JAVA) -Xmx500M -cp "$(BUILD_DIR_FULLPATH)/deps/antlr-$(ANTLR_VERSION)-complete.jar:$$CLASSPATH" org.antlr.v4.Tool -Dlanguage=Go -package dmlex ./dml.g4
	$(file > $@,)

define CHECK_GENERATED_CLEAN
if [ -n "$$(git status --porcelain -- $(G4_DML) ./pkg/dmlex)" ]; then \
	echo "Error: generated files have uncommitted changes. Please run 'make generate' and commit the result."; \
	git status --short -- $(G4_DML) ./pkg/dmlex; \
	git diff -- $(G4_DML) ./pkg/dmlex; \
	exit 1; \
fi
endef

DMLEX_GENERATED_FLAG:=$(BUILD_ENV_DIR)/dmlex-generated.flag
$(DMLEX_GENERATED_FLAG): $(G4_GENERATED_FLAG) $(BUILD_ENV_DIR)
$(DMLEX_GENERATED_FLAG):
	rm -f ./pkg/dmlex/dml_parser*.go
	cp $(filter-out $(G4_DML_DIR)/dml_parser.go,$(G4_GENERATED_FILES)) ./pkg/dmlex/
	(cd tools/baker && $(GO) run . -out $(realpath .)/pkg/dmlex $(realpath .)/$(G4_DML_DIR)/dml_parser.go > /dev/null)
	rm -f $(G4_GENERATED_FILES)
ifneq ($(filter-out 0 false FALSE no NO,$(CI)),)
	@$(CHECK_GENERATED_CLEAN)
endif
	$(file > $@,)

GO_GENERATED_FLAG:=$(BUILD_ENV_DIR)/go-generated.flag
ifeq ($(filter-out 0 false FALSE no NO,$(SKIP_ANTLR)),)
$(GO_GENERATED_FLAG): $(DMLEX_GENERATED_FLAG) $(BUILD_ENV_DIR)
else
$(GO_GENERATED_FLAG): $(BUILD_ENV_DIR)
endif
$(GO_GENERATED_FLAG):
	$(GO) generate ./pkg/dmlex/...
	$(file > $@,)

ifeq ($(filter-out 0 false FALSE no NO,$(SKIP_ANTLR)),)
generate: $(DMLEX_GENERATED_FLAG) $(GO_GENERATED_FLAG)
ifneq ($(filter-out 0 false FALSE no NO,$(CI)),)
	@$(CHECK_GENERATED_CLEAN)
endif
$(DMLEX_GENERATED_FILES): $(DMLEX_GENERATED_FLAG) $(GO_GENERATED_FLAG)
else
generate: $(GO_GENERATED_FLAG)
$(DMLEX_GENERATED_FILES): $(GO_GENERATED_FLAG)
endif

build-gocache: $(DMLEX_GENERATED_FILES)
	$(GO) build ./...

$(BUILD_DIR)/bin/%: build-gocache
	$(GO) build $(GO_BUILD_TAGS) -ldflags "$(GO_LDFLAGS)" -o $@ "./cmd/$*/..."

APPS_BUILD_WILDCARD:=$(foreach _app,$(BUILD_APPS),$(BUILD_DIR)/bin/$(_app))

WASM_DIR := cmd/dml-web-format
WASM_OUT := $(WASM_DIR)/main.wasm
WASM_EXEC_OUT := $(WASM_DIR)/wasm_exec.js
GOROOT_PATH := $(shell $(GO) env GOROOT)

$(WASM_EXEC_OUT):
	@if [ -f "$(GOROOT_PATH)/misc/wasm/wasm_exec.js" ]; then \
		cp "$(GOROOT_PATH)/misc/wasm/wasm_exec.js" $@; \
	else \
		cp "$(GOROOT_PATH)/lib/wasm/wasm_exec.js" $@; \
	fi

build-wasm: build-gocache $(WASM_EXEC_OUT)
	GOOS=js GOARCH=wasm $(GO) build -o $(WASM_OUT) ./$(WASM_DIR)

# Build profiles {{{

release: GO_BUILD_TAGS:=-tags dmltools.ndebug
release: $(APPS_BUILD_WILDCARD) build-wasm

debug: $(APPS_BUILD_WILDCARD) build-wasm

# }}}

web: build-wasm
	open "http://localhost:8080/"
	$(PYTHON) -m http.server --directory cmd/dml-web-format --bind 127.0.0.1 8080

test: build-gocache
	$(GO) test -timeout 60s ./test/... ./cmd/... ./pkg/...

test-ndebug: build-gocache
	$(GO) test -tags dmltools.ndebug -timeout 60s ./test/... ./cmd/... ./pkg/...

coverage: build-gocache
	$(GO) test -timeout 60s -coverprofile=coverage.out -coverpkg=./pkg/...,./cmd/... ./test/... ./cmd/... ./pkg/...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@$(GO) tool cover -func=coverage.out | tail -n 1

coverage.out coverage.html: coverage
	@true

clean-gocache:
	$(GO) clean -cache

clean-go-testcache:
	$(GO) clean -testcache

clean-go-fuzzcache:
	$(GO) clean -fuzzcache

clean-coverage:
	rm -f coverage.out coverage.html

clean: clean-gocache clean-go-testcache clean-go-fuzzcache clean-coverage

goreleaser-check:
	goreleaser check

goreleaser-build:
	goreleaser build --snapshot --clean

goreleaser-snapshot:
	goreleaser release --snapshot --clean

web-dist: goreleaser-snapshot
	$(PYTHON) -m http.server --directory dist/web --bind 127.0.0.1 8080

doc:
	@missing=$$($(GO) list -f '{{if not .Doc}}{{.ImportPath}}{{end}}' ./pkg/...); \
	if [ -n "$$missing" ]; then \
		echo "Error: the following packages are missing package documentation:"; \
		echo "$$missing"; \
		exit 1; \
	fi
	@for pkg in $$($(GO) list ./pkg/...); do \
		echo ""; \
		echo "=== $$pkg ==="; \
		$(GO) doc $$pkg; \
	done

doc-web:
	@if ! command -v $(GODOC) > /dev/null 2>&1 && [ ! -x "$(GODOC)" ]; then \
		echo "godoc not found. Install it with: go install golang.org/x/tools/cmd/godoc@latest"; \
		exit 1; \
	fi
	@open "http://localhost:6060/pkg/github.com/timoschenko/dmltools/" 2>/dev/null || true
	-$(GODOC) -http=:6060

license-check:
	(cd tools/licensecheck && $(GO) run . -root ../..)

test-tools:
	(cd tools/licensecheck && $(GO) test -v ./...)
	(cd tools/baker && $(GO) test -v ./...)

.PHONY: format generate test test-ndebug coverage clean-coverage clean clean-gocache clean-go-testcache clean-go-fuzzcache web goreleaser-check goreleaser-build goreleaser-snapshot web-dist doc doc-web license-check test-tools
