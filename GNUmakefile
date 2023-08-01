TEST?=$$(go list ./...  | grep -v 'vendor')

PKG_NAME=nexus
PKG_ARCH ?= amd64

GOCMD=go
GOBUILD=$(GOCMD) build
GOFLAGS=-mod=vendor
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

NEXUS_HOST=$(shell cd ./scripts && ./detect-docker-env-ip.sh)
MINIO_HOST=$(shell if [ "$(NEXUS_HOST)" = "127.0.0.1" ]; then echo "minio"; else echo "$(NEXUS_HOST)"; fi;)
NEXUS_PORT=$(shell grep -E "(NEXUS_PORT=)" ./scripts/.env | grep -oE "[0-9]+")
MINIKUBE_MOUNT_PID=$(word 1,$(shell ps | grep -v grep | grep 'minikube mount' | grep $(PWD)/scripts))

default: build

start-services:
ifeq (minikube,$(MINIKUBE_ACTIVE_DOCKERD))
ifeq (,$(MINIKUBE_MOUNT_PID))
	minikube mount $(PWD)/scripts:$(PWD)/scripts --uid=200 --gid=200 &
endif
endif
	cd ./scripts && ./start-services.sh && cd -

stop-services:
ifneq (,$(MINIKUBE_MOUNT_PID))
	kill $(MINIKUBE_MOUNT_PID)
endif
	cd ./scripts && ./stop-services.sh && cd -

build: fmtcheck
	go build -v .

linux: fmtcheck
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o terraform.d/plugins/linux_amd64/terraform-provider-nexus -v

test: fmt
	go test $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test -cover $(TESTARGS) -timeout=30s -parallel=4

testacc: fmt
	NEXUS_URL="http://$(NEXUS_HOST):$(NEXUS_PORT)" \
	NEXUS_USERNAME="admin" \
	NEXUS_PASSWORD="admin123" \
	AWS_ACCESS_KEY_ID="minioadmin" \
	AWS_SECRET_ACCESS_KEY="minioadmin" \
	AWS_ENDPOINT="http://$(MINIO_HOST):9000" \
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

docs:
	go generate ./...

.PHONY: build start-services stop-services test testacc fmt fmtcheck lint tools docs
