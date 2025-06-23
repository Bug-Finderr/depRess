package graph

import "maps"

type Node struct {
	ID      string // package@version
	Name    string
	Version string
	Depth   int
}

type Graph struct {
	Nodes     map[string]*Node    // id: Node
	Edges     map[string][]string // id: [dependent ids]
	InDegree  map[string]int      // id: in-degree count
	Conflicts []Conflict          // [conflicting conflicts]
	Missing   []string            // [missing ids]
}

type Conflict struct {
	Package   string
	Requested string
	Existing  string
	Reason    string
}

func NewGraph() *Graph {
	return &Graph{
		Nodes:    make(map[string]*Node),
		Edges:    make(map[string][]string),
		InDegree: make(map[string]int),
	}
}

func (g *Graph) AddNode(id string, node *Node) {
	g.Nodes[id] = node
	if _, exists := g.InDegree[id]; !exists {
		g.InDegree[id] = 0
	}
}

func (g *Graph) AddEdge(from string, to string) {
	for _, existing := range g.Edges[from] {
		if existing == to {
			return
		}
	}
	g.Edges[from] = append(g.Edges[from], to)
	g.InDegree[to]++
}

func (g *Graph) AddConflict(pkg, req, existing, reason string) {
	g.Conflicts = append(g.Conflicts, Conflict{
		Package:   pkg,
		Requested: req,
		Existing:  existing,
		Reason:    reason,
	})
}

func (g *Graph) AddMissing(pkg string) {
	g.Missing = append(g.Missing, pkg)
}

func (g *Graph) GetNodesSortedByDepth() []*Node {
	var nodes []*Node
	for _, node := range g.Nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (g *Graph) DetectCycles() [][]string {
	visited := make(map[string]bool)
	stack := make(map[string]bool)
	var cycles [][]string

	var dfs func(string, []string) bool // remove panic
	dfs = func(node string, path []string) bool {
		if stack[node] {
			start := -1

			for i, p := range path {
				if p == node {
					start = i
					break
				}
			}

			if start != -1 {
				cycles = append(cycles, append(path[start:], node))
			}

			return true
		}

		if visited[node] {
			return false
		}

		visited[node] = true
		stack[node] = true

		for _, nb := range g.Edges[node] {
			if dfs(nb, append(path, node)) {
				return true
			}
		}

		stack[node] = false
		return false
	}

	for id := range g.Nodes {
		if !visited[id] {
			dfs(id, []string{})
		}
	}

	return cycles
}

func (g *Graph) TopoSort() []string {
	inDegreeCopy := maps.Clone(g.InDegree)
	var q []string
	var ans []string

	for id := range g.Nodes {
		if inDegreeCopy[id] == 0 {
			q = append(q, id)
		}
	}

	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		ans = append(ans, cur)

		for _, nb := range g.Edges[cur] {
			inDegreeCopy[nb]--
			if inDegreeCopy[nb] == 0 {
				q = append(q, nb)
			}
		}
	}

	if len(ans) != len(g.Nodes) {
		return nil
	}

	return ans
}
