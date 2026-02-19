.PHONY: help build run run-dev clean vet fmt tidy install-deps update docker-build docker-compose-up docker-compose-down proto proto-generate

APP_NAME = recording-service
CMD_PATH = ./cmd/recording-service
BIN_DIR = bin
PORT = 8096
PROTO_ROOT = pkg/recording_service
PROTO_FILE = recording.proto
GEN_DIR = pkg/gen/recording_service
GO_MODULE = github.com/psds-microservice/recording-service
.DEFAULT_GOAL := help

help:
	@echo "recording-service (gRPC only)"
	@echo "  make build        - Build binary"
	@echo "  make run          - Build and run (api)"
	@echo "  make run-dev      - Run without build (go run ... api)"
	@echo "  make clean vet fmt tidy update"
	@echo "  make docker-build docker-compose-up"
	@echo "  Port: $(PORT)  gRPC: localhost:$(PORT)"

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_PATH)
	@echo "OK: $(BIN_DIR)/$(APP_NAME)"

run: build
	@cd $(BIN_DIR) && ./$(APP_NAME) api

run-dev:
	go run $(CMD_PATH) api

clean:
	rm -rf $(BIN_DIR)
	go clean

vet:
	go vet ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

install-deps:
	go mod download

update:
	@echo "🔄 Updating dependencies..."
	go get -u ./...
	go mod tidy
	go mod vendor
	$(MAKE) proto-generate
	@echo "✅ Dependencies updated"

proto: proto-generate
proto-generate:
	@mkdir -p $(GEN_DIR); PATH="$$(go env GOPATH 2>/dev/null)/bin:$$PATH"; \
	protoc -I $(PROTO_ROOT) --go_out=. --go_opt=module=$(GO_MODULE) --go-grpc_out=. --go-grpc_opt=module=$(GO_MODULE) $(PROTO_ROOT)/$(PROTO_FILE); \
	echo "OK: $(GEN_DIR)"

docker-build:
	docker build -f deployments/Dockerfile -t $(APP_NAME):latest .

docker-compose-up:
	docker compose -f deployments/docker-compose.yml up -d

docker-compose-down:
	docker compose -f deployments/docker-compose.yml down
