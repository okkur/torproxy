BIN := torproxy
MAINTAINER := okkur
VERSION := `cat $(VERSION_FILE)`
IMAGE := $(MAINTAINER)/$(BIN):$(VERSION)

BUILD_GOOS := $(if $(GOOS),$(GOOS),linux)
BUILD_GOARCH := $(if $(GOARCH),$(GOARCH),amd64)

CONTAINER ?= $(BIN)

.DEFAULT_GOAL := build

build:
	cd cmd/torproxy && \
	GO111MODULE=on CGO_ENABLED=0 GOARCH=$(BUILD_GOARCH) GOOS=$(BUILD_GOOS) go build -ldflags="-s -w"
	mv cmd/torproxy/torproxy ./$(BIN)

test:
	GO111MODULE=on go test -v `go list ./...`

image-build: docker-build
	docker build -t $(IMAGE) .

docker-run: image-build
	docker run --name $(CONTAINER) $(IMAGE)

docker-test:
	docker run --network=host -v $(shell pwd):/source -v $(GOPATH)/pkg/mod:/go/pkg/mod golang:1.13-alpine /bin/sh \
	-c "cd /source && apk add tor git gcc musl-dev make && GOROOT=\"/usr/local/go\" make test"

docker-build:
	docker run --network=host -v $(shell pwd):/source -v $(GOPATH)/pkg/mod:/go/pkg/mod golang:1.13-alpine /bin/sh \
	-c "cd /source && apk add tor git gcc musl-dev make && make build"

version:
	@echo $(VERSION)
