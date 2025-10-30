.PHONY: tidy fmt vet lint build run test cover

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

build:
	go build -trimpath -ldflags="-s -w" -o bin/app ./cmd/app

run:
	go run ./cmd/app

test:
	go test ./... -race -shuffle=on -cover

cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
