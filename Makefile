all: test fmt lint run

test:
	@GO111MODULE=on go test -cover -race ./...

build:
	@GO111MODULE=on go build -o service

clean: docker-clean
	@rm service

deps:
	@GO111MODULE=on go mod download
	@GO111MODULE=on go mod tidy

docker: fmt lint build
	@docker build . -t indiependente/barcode

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