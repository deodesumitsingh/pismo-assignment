APP_NAME := pismo_app
BUILD_DIR := bin
SRC_DIR := .
PKG := ./...

.PHONY: all test build dev-run run

all: test build 

test:
	go test -v $(PKG)

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)/main.go

dev-run: build
	./$(BUILD_DIR)/$(APP_NAME)

run: build 
	GIN_MODE=release ./$(BUILD_DIR)/$(APP_NAME)
