build:
	go build -o birjasmm ./cmd/bot

run:
	go run ./cmd/bot

test:
	go test ./...

tidy:
	go mod tidy

.PHONY: build run test tidy
