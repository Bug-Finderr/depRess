BUILD_DIR=./bin

.PHONY: build clean run deps help

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/depRess .

clean:
	rm -rf $(BUILD_DIR)
	go clean

run:
	$(BUILD_DIR)/depRess $(if $(file),-file $(file)) $(if $(depth),-depth $(depth))

deps:
	go mod tidy
	go mod download

help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  clean      - Clean build artifacts"
	@echo "  run        - Build and run with optional file=path and depth=number flags"
	@echo "  deps       - Install dependencies"
	@echo "  help       - Show this help"
