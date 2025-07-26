build:
	@go build -o bin/foorrent main.go

run:
	@go run main.go

clean:
	@rm -rf bin

test:
	@go test -v -failfast ./...

.PHONY: build run clean test
