VERSION := $(shell git rev-parse --short HEAD)

install:
	go install -ldflags '-X github.com/hacker65536/asg/cmd.GitCommit=$(VERSION)'

build: test
	go build -o ~/.local/bin/

test: fmt vet 
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic

vet:
	go vet ./...


fmt:
	go fmt ./...
