package graph

import "maps"

type Node struct {
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

	var dfs func(string, []string) // remove panic
	dfs = func(node string, path []string) {
		if stack[node] {
			for i, p := range path {
				if p == node {
					cycles = append(cycles, append(path[i:], node))
					return
				}
			}
		}

		if visited[node] {
			return
		}

		visited[node] = true
		stack[node] = true

		for _, nb := range g.Edges[node] {
			dfs(nb, append(path, node))
		}

		stack[node] = false
	}

	for id := range g.Nodes {
		if !visited[id] {
			dfs(id, []string{})
		}
	}

	return cycles
}

func (g *Graph) TopoSort() []string {
	inDegree := maps.Clone(g.InDegree)
	var q, ans []string

	for id := range g.Nodes {
		if inDegree[id] == 0 {
			q = append(q, id)
		}
	}

	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		ans = append(ans, cur)

		for _, nb := range g.Edges[cur] {
			inDegree[nb]--
			if inDegree[nb] == 0 {
				q = append(q, nb)
			}
		}
	}

	if len(ans) != len(g.Nodes) {
		return nil
	}

	return ans
}
