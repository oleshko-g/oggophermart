.DEFAULT_GOAL := build

.PHONY: fmt vet test build gen

gen:
	goa gen github.com/oleshko-g/oggophermart/api/design -o internal/
	sqlc generate

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
