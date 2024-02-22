APP_NAME=viego-art-backend
BUILD_DIR=./build
SRC_DIR=./cmd/viego-art-backend

all: build run

build:
	@echo "Building $(APP_NAME)"
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)

run:
	@echo "Running $(APP_NAME)"
	@$(BUILD_DIR)/$(APP_NAME)

clean:
	@echo "Cleaning"
	del /s /q "$(BUILD_DIR)"

.PHONY: all build run clean