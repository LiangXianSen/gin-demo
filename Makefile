GO=go

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

TARGET = demo

project = github.com/LiangXianSen/gin-demo
packages := $(shell go list ./...|grep -v /vendor/)

.PHONY: check test lint

all: build

run:
	@./demo

build: $(TARGET)

$(TARGET): $(SRC)
	@$(GO) build $(project)/cmd/$@

test: check
	@$(GO) test -race $(packages) -v -coverprofile=.coverage.out
	@$(GO) tool cover -func=.coverage.out
	@rm -f .coverage.out

check:
	@$(GO) vet -composites=false $(packages)

lint:
	@golangci-lint run ./...

doc:
	@godoc -http=localhost:8089

clean:
	@rm $(TARGET)
