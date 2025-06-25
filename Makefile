BUILD_DIR=./bin

.PHONY: build clean run deps docker-build docker-run docker-dot-png docker-dot-svg docker-clean help

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

docker-dot-png:
	@docker run --rm --entrypoint dot -v $(PWD):/output depres -Tpng dependency_graph.dot -o dependency_graph.png
	@echo "dependency_graph.png created successfully."

docker-dot-svg:
	@docker run --rm --entrypoint dot -v $(PWD):/output depres -Tsvg dependency_graph.dot -o dependency_graph.svg
	@echo "dependency_graph.svg created successfully."

docker-clean:
	docker rmi depres || true

help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Build and run the application with optional flags:"
	@echo "                   file=<path>   (Path to package.json, default: package.json)"
	@echo "                   depth=<number> (Dependency resolution depth, default: 3)"
	@echo "                   viz=1         (Generate graphviz DOT file)"
	@echo "  deps         - Install dependencies"
	@echo "  docker-build - Build the Docker image"
	@echo "  docker-run   - Run the Docker container with optional flags:"
	@echo "                   file=<path>   (Path to package.json, default: package.json)"
	@echo "                   depth=<number> (Dependency resolution depth, default: 3)"
	@echo "                   viz=1         (Generate graphviz DOT file)"
	@echo "  docker-dot-png - Convert dependency_graph.dot to PNG using Docker"
	@echo "  docker-dot-svg - Convert dependency_graph.dot to SVG using Docker"
	@echo "  docker-clean - Remove the depres Docker image"
	@echo "  help         - Show this help message"
	@echo "Examples:"
	@echo "  make run file=./path/to/package.json depth=2 viz=1"
	@echo "  make docker-run file=package-warning.json depth=5 viz=1"
	@echo "  make docker-run viz=1 && make docker-dot-png"
