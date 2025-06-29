.PHONY: build test generate fmt lint tidy check cover cover-html cover-cli mock run swag
.SILENT: cover

include .env
export

build:
	CGO_ENABLED=0 go build -buildvcs=false -o ./bin/gophermart ./cmd/gophermart

test:
	go test ./... -v

generate:
	go generate ./...
	make format

fmt:
	goimports -l -w .
	swag fmt

lint:
	golangci-lint run

tidy:
	go mod tidy

check: fmt tidy build lint cover-cli

cover:
	go test ./... -coverpkg='./internal/...', -coverprofile coverage-temp.out
	cat coverage-temp.out | grep -v "_easyjson.go\|mocks" > coverage.out
	rm coverage-temp.out

cover-html: cover
	go tool cover -html=coverage.out

cover-cli: cover
	go tool cover -func=coverage.out

mock:
	mockgen -destination=internal/storage/storage_mocks/mock_repository.go -package=storage_mocks github.com/cmrd-a/gophermart/internal/storage Repository

swag:
	swag init --parseDependency --generalInfo server.go --dir internal/api --output internal/api/docs --parseInternal true
	swag fmt

run:build
	./bin/gophermart
	