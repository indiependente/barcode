PROJECT = $(notdir $(shell pwd))
SERVICE_DIRS = $(shell ls services)
THIS_FILE = $(lastword $(MAKEFILE_LIST))
PLATFORM_VERSION = $(shell git rev-parse --short HEAD)

.PHONY: all
all: test fmt lint

.PHONY: test
test:
	@GO111MODULE=on go test -cover -race ./...

.PHONY: build-all
build-all:
	@for dir in $(SERVICE_DIRS); do \
		@GO111MODULE=on go build -ldflags="-s -w" -o $$dir \
	done

.PHONY: clean
clean: docker-clean
	@rm service

.PHONY: deps
deps:
	rm go.mod go.sum
	@GO111MODULE=on go mod init
	@GO111MODULE=on go mod tidy
	@GO111MODULE=on go mod download

.PHONY: docker-build
docker-build:
	@docker build -t indiependente/$(PROJECT)-$(SERVICE) -f services/$(SERVICE)/Dockerfile .

.PHONY: docker-tag
docker-tag:
	@docker tag indiependente/$(PROJECT)-$(SERVICE) indiependente/$(PROJECT)-$(SERVICE):$(VERSION)

.PHONY: docker-build-all 
docker-build-all:
	@for dir in $(SERVICE_DIRS) ; do \
		$(MAKE) -f $(THIS_FILE) docker-build SERVICE=$$dir ;\
		$(MAKE) -f $(THIS_FILE) docker-tag SERVICE=$$dir VERSION=$(PLATFORM_VERSION); \
	done

docker-clean:
	@docker rmi $$(docker images | grep barcode | awk '{print $$3}') --force

fmt:
	@GO111MODULE=on go fmt ./...

lint:
	@command -v golangci-lint || (cd /usr/local ; wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s latest)
	@golangci-lint --version
	@golangci-lint run --disable-all \
	--deadline=10m \
	--skip-dirs vendor \
	--skip-files \.*_mock\.*\.go \
	-E errcheck \
	-E govet \
	-E unused \
	-E gocyclo \
	-E golint \
	-E varcheck \
	-E structcheck \
	-E maligned \
	-E ineffassign \
	-E interfacer \
	-E unconvert \
	-E goconst \
	-E gosimple \
	-E staticcheck \
	-E gosec

run:
	@GO111MODULE=on go run main.go
