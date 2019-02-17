all: test run

test:
	GO111MODULE=on go test ./...

build:
	GO111MODULE=on go build -o service

run:
	GO111MODULE=on go run main.go
