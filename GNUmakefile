TEST?=$$(go list ./... | grep -v 'vendor')
PKG_NAME=nexus

GO111MODULE111=on
GOFLAGS=-mod=vendor

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -w -s ./$(PKG_NAME)

fmtcheck:
	@echo "==> Checking source code againt gofmt..."
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	@golangci-lint run ./$(PKG_NAME)

tools:
	@echo "==> installing required tooling..."
	go install github.com/client9/misspell/cmd/misspell
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: build test fmt fmtcheck lint tools