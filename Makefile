VCS_REF    = $(shell git rev-parse --short HEAD)
VERSION    = v$(shell git describe --always --match "v*")
TAG        = rafaeljusto/teamwork-ai:$(VERSION)
LATEST_TAG = rafaeljusto/teamwork-ai:latest

.PHONY: build deploy artifacts

default: build

build:
	docker build .

deploy:
	docker buildx build \
	  --platform linux/amd64,linux/arm64 \
		--build-arg BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
	  --build-arg BUILD_VCS_REF=$(shell git rev-parse --short HEAD) \
	  --build-arg BUILD_VERSION=$(VERSION) \
	  -t $(TAG) \
	  -t $(LATEST_TAG) \
	  --push \
	  --progress=plain \
	  .

artifacts:
	GOOS=windows GOARCH=amd64 go build -o teamwork-mcp-windows-amd64 ./cmd/mcp
	GOOS=linux   GOARCH=amd64 go build -o teamwork-mcp-linux-amd64 ./cmd/mcp
	GOOS=linux   GOARCH=arm64 go build -o teamwork-mcp-linux-arm64 ./cmd/mcp
	GOOS=darwin  GOARCH=amd64 go build -o teamwork-mcp-darwin-amd64 ./cmd/mcp
	GOOS=darwin  GOARCH=arm64 go build -o teamwork-mcp-darwin-arm64 ./cmd/mcp