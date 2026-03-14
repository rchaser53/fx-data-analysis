BIN_DIR ?= bin

.PHONY: build-fetch-usdjpy up clean

up:
	./start-app.sh

build-fetch-usdjpy:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o $(BIN_DIR)/fetch-usdjpy ./cmd/fetch-usdjpy

clean:
	@rm -rf $(BIN_DIR)
