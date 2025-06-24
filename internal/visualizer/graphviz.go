package visualizer

import (
	"depRess/internal/graph"
	"fmt"
	"os"
	"strings"
)

type Graphviz struct {
	graph *graph.Graph
}

func New(g *graph.Graph) *Graphviz {
	return &Graphviz{graph: g}
}

func (v *Graphviz) Generate() error {
	var dot strings.Builder

	dot.WriteString("digraph dependencies {\n")
	dot.WriteString("  rankdir=TB;\n")
	dot.WriteString("  node [shape=box, style=filled, fontname=\"Arial\"];\n")
	dot.WriteString("  edge [fontname=\"Arial\", fontsize=10];\n\n")

	for _, node := range v.graph.Nodes {
		color := v.getNodeColor(node.Depth)
		fillColor := v.getNodeFillColor(node.Depth)
		nodeName := v.cleanNodeName(node.Name)
		label := fmt.Sprintf("%s\\n%s", node.Name, node.Version)

		dot.WriteString(fmt.Sprintf("  \"%s\" [label=\"%s\", color=\"%s\", fillcolor=\"%s\"];\n", nodeName, label, color, fillColor))
	}

	dot.WriteString("\n  // Edges\n")

	// add edges using adj list
	for from, edges := range v.graph.Edges {
		fromName := v.cleanNodeName(from)
		for _, to := range edges {
			toName := v.cleanNodeName(to)
			dot.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\";\n", fromName, toName))
		}
	}

	// conflicts as red dashed edges
	if len(v.graph.Conflicts) > 0 {
		dot.WriteString("\n  // Conflicts\n")
		for _, conflict := range v.graph.Conflicts {
			conflictName := v.cleanNodeName(conflict.Package)
			dot.WriteString(fmt.Sprintf("  \"%s\" [color=red, style=\"filled,dashed\"];\n", conflictName))
		}
	}

	// missing pkgs as gray nodes
	if len(v.graph.Missing) > 0 {
		dot.WriteString("\n  // Missing packages\n")
		for _, missing := range v.graph.Missing {
			missingName := v.cleanNodeName(missing)
			dot.WriteString(fmt.Sprintf("  \"%s\" [label=\"%s\\n(MISSING)\", color=gray, fillcolor=lightgray, style=\"filled,dashed\"];\n",
				missingName, missing))
		}
	}

	dot.WriteString("}\n")
	return os.WriteFile("dependency_graph.dot", []byte(dot.String()), 0644)
}

func (v *Graphviz) getNodeColor(depth int) string {
	colors := []string{"blue", "green", "yellow", "orange", "red"}
	if depth < 0 || depth >= len(colors) {
		return "black"
	}
	return colors[depth]
}

func (v *Graphviz) getNodeFillColor(depth int) string {
	colors := []string{"lightblue", "lightgreen", "lightyellow", "lightpink", "lightcyan"}
	if depth < 0 || depth >= len(colors) {
		return "white"
	}
	return colors[depth]
}

func (v *Graphviz) cleanNodeName(name string) string {
	// Replace problematic characters for DOT format
	name = strings.ReplaceAll(name, "@", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")
	return name
}
