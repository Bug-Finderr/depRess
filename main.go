package main

import (
	"depRess/internal/resolver"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	filePath := flag.String("file", "package.json", "Path to package.json file")
	maxDepth := flag.Int("depth", 3, "Maximum dependency resolution depth (1-5)")
	viz := flag.Bool("viz", false, "Generate graphviz DOT file for visualization")
	help := flag.Bool("help", false, "Show help message")
	flag.Parse()

	if *help {
		demo()
		return
	}

	if *maxDepth < 1 || *maxDepth > 5 {
		fmt.Println("Error: max-depth must be between 1 and 5")
		os.Exit(1)
	}

	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		fmt.Printf("Error: package.json file not found at %s\n", *filePath)
		os.Exit(1)
	}

	fmt.Println("Node Dependency Resolver Starting...")
	fmt.Printf("Package.json: %s\n", *filePath)
	fmt.Printf("Max Depth: %d\n", *maxDepth)
	fmt.Println(strings.Repeat("-", 60))

	res := resolver.New(*maxDepth)

	if err := res.ResolveDeps(*filePath); err != nil {
		fmt.Printf("Error during resolution: %v\n", err)
		os.Exit(1)
	}

	res.GenReport()

	if *viz {
		res.GenViz()
	}
}

func demo() {
	fmt.Println("Node Dependency Resolver - by Bug")
	fmt.Println("\nUsage:")
	fmt.Printf("  %s [flags]\n", os.Args[0])
	fmt.Println("\nFlags:")
	fmt.Println("\t-file string")
	fmt.Println("\t\t\tPath to package.json file (default \"package.json\")")
	fmt.Println("\t-depth int")
	fmt.Println("\t\t\tMaximum dependency resolution depth (1-5) (default 3)")
	fmt.Println("\t-viz")
	fmt.Println("\t\t\tGenerate graphviz DOT file for visualization")
	fmt.Println("\t-help")
	fmt.Println("\t\t\tShow this help message")
	fmt.Println("\nExamples:")
	fmt.Printf("  %s\n", os.Args[0])
	fmt.Printf("  %s -file ./<path>/package.json\n", os.Args[0])
	fmt.Printf("  %s -depth 2 -file package.json -viz\n", os.Args[0])
	fmt.Println("\nVisualization:")
	fmt.Println("\tInstall: brew install graphviz")
	fmt.Println("\tRender: dot -Tpng dependency_graph.dot -o graph.png")
	fmt.Println("\tMore info: https://graphviz.org/doc/info/command.html")
}
