.PHONY: build test generate format lint tidy check cover cover-html cover-cli mock run si

build:
	go build -o ./bin/gophermart ./cmd/gophermart

test:
	go test ./... -v

generate:
	go generate ./...
	make format

format:
	goimports -l -w .
	swag fmt

lint:
	go vet ./...
	staticcheck -checks=all,-ST1000, ./...

tidy:
	go mod tidy

check: build tidy format lint test

cover:
	go test ./... -coverpkg='./internal/...', -coverprofile coverage.out.tmp
	cat coverage.out.tmp | grep -v "_easyjson.go\|mocks" > coverage.out
	rm coverage.out.tmp

cover-html: cover
	go tool cover -html=coverage.out

cover-cli: cover
	go tool cover -func=coverage.out

mock:
	mockgen -destination=internal/storage/storage_mocks/mock_repository.go -package=storage_mocks github.com/cmrd-a/gophermart/internal/storage Repository

si:
	swag init --generalInfo internal/api/server.go --output internal/api/docs

run:build
	./bin/gophermart
	