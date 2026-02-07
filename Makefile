BINARY := blueprint
PKG := github.com/ale/blueprint
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -s -w \
	-X $(PKG)/internal/version.Version=$(VERSION) \
	-X $(PKG)/internal/version.Commit=$(COMMIT) \
	-X $(PKG)/internal/version.Date=$(DATE)

.PHONY: build test clean lint run status

## build: Compila o binario em bin/
build:
	go build -ldflags '$(LDFLAGS)' -o bin/$(BINARY) ./cmd/blueprint

## test: Roda todos os testes
test:
	go test ./... -v

## clean: Remove artefatos de build
clean:
	rm -rf bin/

## lint: Roda go vet
lint:
	go vet ./...

## run: Resolve dependencias, compila e abre o TUI
run:
	@bash scripts/install.sh

## status: Mostra estado detalhado dos modulos
status: build
	@./bin/$(BINARY) status
