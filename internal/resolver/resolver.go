package resolver

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"depRess/internal/graph"
	"depRess/internal/registry"
	"depRess/internal/version"
)

type InputFile struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

type QueueItem struct {
	Name       string
	Constraint string
	Depth      int
}

type Resolver struct {
	maxDepth       int
	graph          *graph.Graph
	resolved       map[string]string // pkg -> resolved_ver
	queue          []QueueItem
	registryClient *registry.Client
}

func New(maxDepth int) *Resolver {
	return &Resolver{
		maxDepth:       maxDepth,
		graph:          graph.NewGraph(),
		resolved:       make(map[string]string),
		registryClient: registry.New(),
	}
}

func (r *Resolver) ParseInputFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading package.json: %w", err)
	}

	var pkg InputFile
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, fmt.Errorf("error parsing package.json: %w", err)
	}

	deps := make(map[string]string)

	for name, ver := range pkg.Dependencies {
		deps[name] = ver
	}

	for name, ver := range pkg.DevDependencies {
		deps[name] = ver
	}

	return deps, nil
}

func (r *Resolver) ResolveDeps(path string) error {
	deps, err := r.ParseInputFile(path)
	if err != nil {
		return err
	}

	fmt.Printf("Found %d dependencies to resolve\n", len(deps))

	// init q with root deps
	for n, c := range deps {
		r.queue = append(r.queue, QueueItem{
			Name:       n,
			Constraint: c,
			Depth:      0,
		})
	}

	processed := 0
	for len(r.queue) > 0 {
		item := r.queue[0]
		r.queue = r.queue[1:]

		// skip if already resolved or max depth reached
		if _, exists := r.resolved[item.Name]; exists || item.Depth >= r.maxDepth {
			continue
		}

		processed++
		if processed%10 == 0 {
			fmt.Printf("Processed %d packages...\n", processed)
		}

		if err := r.resolveSinglePackage(item.Name, item.Constraint, item.Depth); err != nil {
			fmt.Printf("Error resolving %s: %v\n", item.Name, err)
		}
	}

	r.buildAdjList()

	fmt.Println("Dependency resolution complete!")
	return nil
}

func (r *Resolver) resolveSinglePackage(pkg, verConstraint string, depth int) error {
	pkgInfo, err := r.registryClient.GetPkgInfo(pkg)
	if err != nil {
		fmt.Printf("%sPackage not found: %s\n", strings.Repeat("\t", depth), pkg)
		r.graph.AddMissing(pkg)
		return err
	}

	var availableVer []string
	for v := range pkgInfo.Versions {
		availableVer = append(availableVer, v)
	}

	if len(availableVer) == 0 {
		fmt.Printf("%sNo versions available for %s\n", strings.Repeat("\t", depth), pkg)
		r.graph.AddMissing(pkg)
		return fmt.Errorf("no versions available for %s", pkg)
	}

	bestVer, err := version.FindBestVersion(availableVer, verConstraint)
	if err != nil || bestVer == "" {
		fmt.Printf("%sVersion conflict: %s@%s\n", strings.Repeat("\t", depth), pkg, verConstraint)
		r.graph.AddConflict(pkg, verConstraint, "", "No matching version")
		return fmt.Errorf("no matching version for %s@%s", pkg, verConstraint)
	}

	// check for version conflicts with already resolved packages
	if existingVer, exists := r.resolved[pkg]; exists {
		if existingVer != bestVer {
			fmt.Printf("%sVersion conflict: %s needs %s but %s already resolved\n",
				strings.Repeat("\t", depth), pkg, bestVer, existingVer)
			r.graph.AddConflict(pkg, bestVer, existingVer, "Version mismatch")
			return nil
		}
		return nil // already resolved with same ver, skip
	}

	r.resolved[pkg] = bestVer

	node := &graph.Node{ // add to graph
		ID:      pkg, // use name as ID
		Name:    pkg,
		Version: bestVer,
		Depth:   depth,
	}
	r.graph.AddNode(pkg, node)

	// process deps of this package
	if depth < r.maxDepth-1 {
		if versionInfo, exists := pkgInfo.Versions[bestVer]; exists {
			for depName, depConstraint := range versionInfo.Dependencies {
				r.queue = append(r.queue, QueueItem{ // schedule for resolution
					Name:       depName,
					Constraint: depConstraint,
					Depth:      depth + 1,
				})
			}
		}
	}

	return nil
}

func (r *Resolver) buildAdjList() {
	for pkg, ver := range r.resolved {
		pkgInfo, err := r.registryClient.GetPkgInfo(pkg)
		if err != nil {
			continue
		}
		if versionInfo, exists := pkgInfo.Versions[ver]; exists {
			for depName := range versionInfo.Dependencies {
				if _, depResolved := r.resolved[depName]; depResolved {
					r.graph.AddEdge(pkg, depName)
				}
			}
		}
	}
}
