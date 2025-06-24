BUILD_DIR=./bin

.PHONY: build clean run deps help

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/depRess .

clean:
	rm -rf $(BUILD_DIR)
	go clean

run:
	$(BUILD_DIR)/depRess $(if $(file),-file $(file)) $(if $(depth),-depth $(depth)) $(if $(viz),-viz)

deps:
	go mod tidy
	go mod download

docker-build:
	docker build -t depres .

docker-run:
	docker run --rm -v $(PWD):/output depres $(if $(file),-file $(file)) $(if $(depth),-depth $(depth)) $(if $(viz),-viz)

help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  clean      - Clean build artifacts"
	@echo "  run        - Build and run the application with optional flags:"
	@echo "                file=<path>  (Path to package.json, default: package.json)"
	@echo "                depth=<number> (Dependency resolution depth, default: 3)"
	@echo "                viz          (Generate graphviz DOT file)"
	@echo "  deps       - Install dependencies"
	@echo "  help       - Show this help message"
	@echo "Examples:"
	@echo "  make run file=./path/to/package.json depth=2 viz"
