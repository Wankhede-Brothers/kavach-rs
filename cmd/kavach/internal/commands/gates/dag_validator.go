// Package gates provides hook gates for Claude Code.
// dag_validator.go: DAG cycle detection algorithm.
// DACE: Micro-modular split from dag.go
package gates

// DAGValidator validates agent delegation chains for cycles.
type DAGValidator struct {
	edges  map[string][]string
	colors map[string]int // 0=white, 1=gray, 2=black
}

const (
	colorWhite = 0 // Unvisited
	colorGray  = 1 // In current DFS path
	colorBlack = 2 // Fully processed
)

// NewDAGValidator creates a new DAG validator.
func NewDAGValidator() *DAGValidator {
	return &DAGValidator{
		edges:  make(map[string][]string),
		colors: make(map[string]int),
	}
}

// AddEdge adds a delegation edge from parent to child agent.
func (d *DAGValidator) AddEdge(parent, child string) {
	d.edges[parent] = append(d.edges[parent], child)
	if _, ok := d.colors[parent]; !ok {
		d.colors[parent] = colorWhite
	}
	if _, ok := d.colors[child]; !ok {
		d.colors[child] = colorWhite
	}
}

// DetectCycle checks if the graph contains a cycle.
func (d *DAGValidator) DetectCycle() []string {
	for k := range d.colors {
		d.colors[k] = colorWhite
	}

	for node := range d.colors {
		if d.colors[node] == colorWhite {
			if path := d.dfs(node, []string{}); path != nil {
				return path
			}
		}
	}
	return nil
}

func (d *DAGValidator) dfs(node string, path []string) []string {
	d.colors[node] = colorGray
	path = append(path, node)

	for _, neighbor := range d.edges[node] {
		if d.colors[neighbor] == colorGray {
			for i, n := range path {
				if n == neighbor {
					return append(path[i:], neighbor)
				}
			}
			return append(path, neighbor)
		}
		if d.colors[neighbor] == colorWhite {
			if cyclePath := d.dfs(neighbor, path); cyclePath != nil {
				return cyclePath
			}
		}
	}

	d.colors[node] = colorBlack
	return nil
}

// ValidatePath checks if a path contains a cycle (simple version).
func ValidatePath(path []string) (bool, []string) {
	seen := make(map[string]int)
	for i, node := range path {
		if prevIdx, exists := seen[node]; exists {
			return false, path[prevIdx:]
		}
		seen[node] = i
	}
	return true, nil
}
