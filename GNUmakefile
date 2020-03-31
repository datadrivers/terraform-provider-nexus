TEST?=$$(go list ./...  | grep -v 'vendor')

PKG_NAME=nexus
PKG_OS ?= darwin linux
PKG_ARCH ?= amd64

GOCMD=go
GOBUILD=$(GOCMD) build
GO111MODULE111=on
GOFLAGS=-mod=vendor
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

WEBSITE_REPO=github.com/hashicorp/terraform-website

default: build

build: fmtcheck
	go build -v .

linux: fmtcheck
	GCO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o terraform.d/plugins/linux_amd64/terraform-provider-nexus -v

darwin: fmtcheck
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o terraform.d/plugins/darwin_amd64/terraform-provider-nexus -v

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test -cover $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -cover -timeout 120m -parallel=4

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	@golangci-lint run ./$(PKG_NAME)

tools:
	@echo "==> installing required tooling..."
	go install github.com/client9/misspell/cmd/misspell
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build test testacc fmt fmtcheck lint tools website website-test