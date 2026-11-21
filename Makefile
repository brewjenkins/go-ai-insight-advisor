APP_NAME := ai-insight-advisor
BIN_DIR  := bin

.PHONY: help run http build test vet fmt tidy clean

help:
	@echo "Available targets:"
	@echo "  make run    - Run HTTP server (listener)"
	@echo "  make build  - Build binary to $(BIN_DIR)/$(APP_NAME)"
	@echo "  make test   - Run unit tests with race detector"
	@echo "  make vet    - Run go vet"
	@echo "  make fmt    - Format Go files with gofmt"
	@echo "  make tidy   - Run go mod tidy"
	@echo "  make clean  - Remove build artifacts"

run:
	go run . listener

http: run

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) .

test:
	go test -race -v ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

tidy:
	go mod tidy

clean:
	rm -rf $(BIN_DIR)
