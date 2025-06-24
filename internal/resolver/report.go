package resolver

import (
	"fmt"
	"strings"
)

func (r *Resolver) GenReport() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("DEPENDENCY RESOLUTION REPORT")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\nSUMMARY:")
	fmt.Printf("   Total packages resolved: %d\n", len(r.resolved))
	fmt.Printf("   Missing packages: %d\n", len(r.graph.Missing))
	fmt.Printf("   Version conflicts: %d\n", len(r.graph.Conflicts))

	if len(r.graph.Missing) > 0 {
		fmt.Printf("\nMISSING PACKAGES (%d):\n", len(r.graph.Missing))
		for _, pkg := range r.graph.Missing {
			fmt.Printf("   • %s\n", pkg)
		}
	}

	if len(r.graph.Conflicts) > 0 {
		fmt.Printf("\nVERSION CONFLICTS (%d):\n", len(r.graph.Conflicts))
		for _, conflict := range r.graph.Conflicts {
			fmt.Printf("   • %s: requested %s, conflict: %s (%s)\n",
				conflict.Package, conflict.Requested, conflict.Existing, conflict.Reason)
		}
	}

	cycles := r.graph.DetectCycles()
	if len(cycles) > 0 {
		fmt.Printf("\nCIRCULAR DEPENDENCIES (%d):\n", len(cycles))
		for i, cycle := range cycles {
			fmt.Printf("   %d. %s\n", i+1, strings.Join(cycle, " -> "))
		}
	} else {
		fmt.Println("\nNO CIRCULAR DEPENDENCIES DETECTED")
	}

	fmt.Println("\nINSTALLATION ORDER (Topological Sort):")
	sorted := r.graph.TopoSort()

	if len(sorted) == len(r.graph.Nodes) {
		fmt.Println("   Valid installation order found:")
		for i := len(sorted) - 1; i >= 0; i-- {
			pkgId := sorted[i]
			if node, exists := r.graph.Nodes[pkgId]; exists {
				depthIndicator := strings.Repeat("  ", node.Depth)
				fmt.Printf("   %2d. %s%s@%s\n", len(sorted)-i, depthIndicator, node.Name, node.Version)
			}
		}
	} else {
		fmt.Println("   Cannot determine installation order due to circular dependencies")
	}

	fmt.Println(strings.Repeat("=", 60))
}
