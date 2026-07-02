BINARY   := adex
MODULE   := github.com/gmvstudio/adex-cli
VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
DATE     := $(shell date +%Y-%m-%d)
LDFLAGS  := -s -w -X $(MODULE)/internal/build.Version=$(VERSION) -X $(MODULE)/internal/build.Date=$(DATE)
PREFIX   ?= /usr/local

.PHONY: all build vet fmt-check test unit-test lint tidy-check install uninstall clean npm-publish oss-sync

all: test

build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY) .

vet:
	go vet ./...

fmt-check:
	@unformatted=$$(gofmt -l . | grep -v '^\.idea/' || true); \
	if [ -n "$$unformatted" ]; then \
		echo "Unformatted Go files:"; \
		echo "$$unformatted"; \
		echo "Run 'gofmt -w .' and commit."; \
		exit 1; \
	fi

# unit-test runs the suite with the race detector, matching the CI gate.
unit-test:
	go test -race -count=1 ./...

# tidy-check fails if go.mod/go.sum are not tidy.
tidy-check:
	go mod tidy
	git diff --exit-code go.mod go.sum

# lint runs golangci-lint pinned to the CI version.
lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6 run

test: vet fmt-check unit-test

install: build
	install -d $(PREFIX)/bin
	install -m755 $(BINARY) $(PREFIX)/bin/$(BINARY)
	@echo "OK: $(PREFIX)/bin/$(BINARY) ($(VERSION))"

uninstall:
	rm -f $(PREFIX)/bin/$(BINARY)

clean:
	rm -f $(BINARY)

npm-publish:
	npm publish --access public

oss-sync:
	./scripts/oss-sync.sh
