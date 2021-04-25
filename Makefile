.PHONY:frontend
GO           ?= go
GOFMT        ?= $(GO)fmt
FIRST_GOPATH := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))

pkgs          = $(shell $(GO) list ./... | grep -v /vendor/)

PREFIX       ?= $(shell pwd)
DIRNAME      ?= $(shell dirname $(shell pwd))

#TAG          ?= $(shell date +%s)
TAG          ?= $(shell git rev-parse --short HEAD)

ENV	     ?= prod

style:
	@echo ">> checking code style"
	@! $(GOFMT) -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)

build:
	@echo ">> go build ..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags -w -o ping-prober-config-generator main.go
	@echo ">> completed."

image:
	@echo ">> go build ..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags -w -o ping-prober-config-generator main.go
	@echo ">> start building ping-prober-config-generator image ..."
	@docker build -t zbd20/ping-prober-config-generator:${VERSION} .
	@echo ">> docker image has been built."
	@echo ">> push image to the repository ..."
	@docker push zbd20/ping-prober-config-generator:${VERSION}
	@echo ">> completed."

clean:
	@echo ">> remove ping-prober-config-generator"
	@rm ping-prober-config-generator