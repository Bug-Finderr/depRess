package main

import (
	"depRess/internal/resolver"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	filePath := flag.String("f", "package.json", "Path to package.json file")
	maxDepth := flag.Int("d", 3, "Maximum dependency resolution depth (1-5)")
	help := flag.Bool("h", false, "Show help message")
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
}

func demo() {
	fmt.Println("Node Dependency Resolver - by Bug")
	fmt.Println("\nUsage:")
	fmt.Printf("  %s [flags]\n", os.Args[0])
	fmt.Println("\nFlags:")
	fmt.Println("\t-f string")
	fmt.Println("\t\t\tPath to package.json file (default \"package.json\")")
	fmt.Println("\t-d int")
	fmt.Println("\t\t\tMaximum dependency resolution depth (1-5) (default 3)")
	fmt.Println("\t-h")
	fmt.Println("\t\t\tShow this help message")
	fmt.Println("\nExamples:")
	fmt.Printf("  %s\n", os.Args[0])
	fmt.Printf("  %s -f ./<path>/package.json\n", os.Args[0])
	fmt.Printf("  %s -d 2 -f package.json\n", os.Args[0])
}
