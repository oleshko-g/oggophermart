.DEFAULT_GOAL := build

.PHONY: fmt vet test build oggoa

oggoa:
	goa gen github.com/oleshko-g/oggophermart/api/design -o internal/

test: vet
	go fmt
	go vet
	revive ./...
	go generate ./...
	go test ./...

build:
	go build

fullbuild: test
	go build
