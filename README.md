# depRess - Node.js Dependency Resolver & Visualizer

[![Go Version](https://img.shields.io/badge/Go-1.24.4-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A command-line tool written in Go that analyzes Node.js `package.json` files, resolves dependency trees with version constraints, detects conflicts, and generates beautiful dependency graphs using graph algorithms.

## 🎥 Demo Video

[🎥 Watch Demo Video](https://drive.google.com/file/d/1WIkfQH6f_OT-MZrpKSOcCIq8kxnESaeh/view?usp=sharing)

## 🎯 Project Overview

**depRess** is a dependency analysis tool (made as a fun project) that combines graph theory algorithms with npm registry data to provide insights into your Node.js project dependencies. It performs deep dependency resolution, conflict detection, and generates visual representations of your dependency tree.

### Core Features

- 🔍 **Dependency Resolution**: Recursively resolves dependencies up to 5 levels deep
- 📊 **Graph-based Analysis**: Uses directed graphs and topological sorting algorithms
- ⚠️ **Conflict Detection**: Identifies version conflicts and missing dependencies
- 📈 **Visual Reports**: Generates detailed text reports with statistics
- 🎨 **Graphviz Visualization**: Creates beautiful dependency graphs as PNG/SVG images
- 🐳 **Docker Support**: Containerized execution for consistent environments
- 🚀 **Great Performance**: Efficient graph algorithms with O(V+E) complexity

## 🏗️ Architecture & Flow

```text
package.json → Parse → NPM Registry API → Version Resolution → Graph Construction → Analysis → Report/Visualization
```

### Core Graph Algorithms

1. **Dependency Resolution**: Modified BFS traversal with version constraint satisfaction
2. **Conflict Detection**: Graph coloring algorithm to identify version conflicts
3. **Topological Sorting**: Kahn's algorithm for dependency ordering
4. **Cycle Detection**: DFS-based cycle detection for circular dependencies
5. **Connected Components**: Identifies isolated dependency clusters

### Project Structure

```text
depRess/
├── main.go                 # CLI entry point and argument parsing
├── internal/
│   ├── graph/             # Graph data structures and algorithms
│   │   └── graph.go       # Node, Edge, Conflict detection
│   ├── registry/          # NPM registry API client
│   │   └── npm.go         # Version fetching and resolution
│   ├── resolver/          # Core dependency resolution engine
│   │   ├── resolver.go    # Main resolution algorithm
│   │   └── report.go      # Report generation
│   ├── version/           # Semantic version handling
│   │   └── resolver.go    # SemVer constraint resolution
│   └── visualizer/        # Graph visualization
│       └── graphviz.go    # DOT file generation
└── public/               # Visualization outputs for examples
    ├── package-success.png
    ├── package-warning.png
    └── package-error.png
```

## 🚀 Installation & Usage

### Prerequisites

- Go 1.24.4 or later
- Docker (Recommended for visualization)
- Graphviz (optional, for local visualization without Docker): `brew install graphviz`

### Build from Source

```bash
# Clone the repository
git clone https://github.com/Bug-Finderr/depRess.git
cd depRess

# Install dependencies
make deps

# Build the binary
make build
```

## 📋 CLI Commands

### Makefile Commands

```bash
# Build project
make build

# Clean build artifacts
make clean

# Install/update dependencies
make deps

# Run with parameters
make run file=package.json depth=3 viz=true

# Docker commands
make docker-build
make docker-run file=package.json depth=2
make docker-dot-png
make docker-dot-svg
make docker-clean

# Show help
make help
```

### Binary Execution (after build)

```bash
# Basic usage
./bin/depRess

# Specify package.json file
./bin/depRess -file ./path/to/package.json

# Set resolution depth (1-5)
./bin/depRess -depth 2

# Generate visualization
./bin/depRess -viz

# Combined flags
./bin/depRess -file package.json -depth 3 -viz

# Show help
./bin/depRess -help
```

### Direct Go Execution

```bash
# Basic run
go run main.go

# With parameters
go run main.go -file package.json -depth 3 -viz

# Help
go run main.go -help
```

### Docker Commands

```bash
# Build Docker image
docker build -t depres .

# Run analysis (mounts current directory)
docker run --rm -v $(PWD):/output depres -file package.json -depth 3

# Run with visualization
docker run --rm -v $(PWD):/output depres -file package.json -viz

# Interactive shell
docker run -it --rm -v $(PWD):/output depres /bin/sh
```

## 🎨 Generating Visualizations

After running the tool with the `-viz` flag, a `dependency_graph.dot` file is created. You can convert this file into an image format like PNG or SVG or PDF using one of the methods below.

### Using Docker (Recommended)

This method leverages the Graphviz installation inside the Docker container, so you don't need to install it on your host machine.

1.  **Generate the `.dot` file using Docker:**

    ```bash
    # This command runs the tool and creates dependency_graph.dot in your project directory
    make docker-run viz=1
    ```

2.  **Convert the `.dot` file to an image:**

    ```bash
    # Generate a PNG image
    make docker-dot-png

    # Or generate an SVG image
    make docker-dot-svg
    ```

### Using a Local Graphviz Installation

If you have Graphviz installed locally, you can use the `dot` command directly:

```bash
# Generate PNG image
dot -Tpng dependency_graph.dot -o dependency_graph.png

# Generate SVG (scalable)
dot -Tsvg dependency_graph.dot -o dependency_graph.svg

# Generate PDF
dot -Tpdf dependency_graph.dot -o dependency_graph.pdf

# Interactive graph (if X11 available)
dot -Tx11 dependency_graph.dot

# Different layouts
dot -Kneato -Tpng dependency_graph.dot -o graph_neato.png
dot -Kcirco -Tpng dependency_graph.dot -o graph_circo.png
dot -Kfdp -Tpng dependency_graph.dot -o graph_fdp.png
```

## 📊 Output Examples

### Text Report

```text
Node Dependency Resolver Starting...
Package.json: package-warning.json
Max Depth: 5
------------------------------------------------------------
Found 4 dependencies to resolve
Version conflict: typescript@5.0.0
Error resolving typescript: no matching version for typescript@5.0.0
Dependency resolution complete!
============================================================
DEPENDENCY RESOLUTION REPORT
============================================================

SUMMARY:
   Total packages resolved: 5
   Missing packages: 0
   Version conflicts: 1

VERSION CONFLICTS (1):
   • typescript: requested 5.0.0, conflict:  (No matching version)

NO CIRCULAR DEPENDENCIES DETECTED

INSTALLATION ORDER (Topological Sort):
   Valid installation order found:
    1.   @radix-ui/react-compose-refs@1.1.2
    2.   scheduler@0.26.0
    3. @radix-ui/react-slot@1.2.3
    4. react-dom@19.1.0
    5. react@19.1.0
============================================================
Graphviz .dot file generated successfully!
```

## 🔧 Configuration

### Command Line Flags

| Flag     | Type   | Default        | Description                     |
| -------- | ------ | -------------- | ------------------------------- |
| `-file`  | string | `package.json` | Path to package.json file       |
| `-depth` | int    | `3`            | Maximum resolution depth (1-5)  |
| `-viz`   | bool   | `false`        | Generate graphviz visualization |
| `-help`  | bool   | `false`        | Show help message               |

## 🧮 Algorithm Complexity

- **Time Complexity**: O(V + E) where V = packages, E = dependency relationships
- **Space Complexity**: O(V + E) for graph storage
- **Network Calls**: O(V) registry API requests with caching

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**@Bug-Finderr** - _Initial work_

---

> _This README was generated by Github Copilot_
