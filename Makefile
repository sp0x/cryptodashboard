APP_NAME = "cryptodashboard"

all: build

build:
	@echo "Building..."
	@go build -o bin/$(APP_NAME) cmd/$(APP_NAME)/main.go
	@echo "Build complete"

serve: build
	@echo "Serving..."
	@./bin/$(APP_NAME)