package glpg

// GLPGProperty represents a map of key-value pairs for properties on nodes and edges.
// It allows for mixed property types.
type GLPGProperty map[string]interface{}

// GLPGNode represents a node (or vertex) in the Generalized Labeled Property Graph.
type GLPGNode struct {
	ID         string       // Unique identifier for the node
	Labels     []string     // List of labels categorizing the node
	Properties GLPGProperty // Key-value pairs storing attributes of the node
}

// GLPGEdge represents an edge (or relationship) in the Generalized Labeled Property Graph.
type GLPGEdge struct {
	ID         string       // Unique identifier for the edge
	SourceID   string       // ID of the source GLPGNode
	TargetID   string       // ID of the target GLPGNode
	Label      string       // Label describing the nature of the relationship
	Properties GLPGProperty // Key-value pairs storing attributes of the edge
}

// GLPG represents the entire Generalized Labeled Property Graph.
type GLPG struct {
	Nodes map[string]*GLPGNode // Map from Node ID to Node
	Edges map[string]*GLPGEdge // Map from Edge ID to Edge

	// Optional: Adjacency lists for easier traversal and performance.
	// These can be populated when the graph is built or queried.
	OutgoingEdges map[string][]*GLPGEdge // Node ID to its outgoing edges
	IncomingEdges map[string][]*GLPGEdge // Node ID to its incoming edges
}

// NewGLPG creates and initializes a new GLPG structure.
func NewGLPG() *GLPG {
	return &GLPG{
		Nodes:         make(map[string]*GLPGNode),
		Edges:         make(map[string]*GLPGEdge),
		OutgoingEdges: make(map[string][]*GLPGEdge),
		IncomingEdges: make(map[string][]*GLPGEdge),
	}
}

// AddNode adds a node to the graph.
// It ensures the adjacency lists are initialized for the node.
func (g *GLPG) AddNode(node *GLPGNode) {
	if node == nil {
		return // Or handle error
	}
	g.Nodes[node.ID] = node
	if _, exists := g.OutgoingEdges[node.ID]; !exists {
		g.OutgoingEdges[node.ID] = []*GLPGEdge{}
	}
	if _, exists := g.IncomingEdges[node.ID]; !exists {
		g.IncomingEdges[node.ID] = []*GLPGEdge{}
	}
}

// AddEdge adds an edge to the graph and updates adjacency lists.
func (g *GLPG) AddEdge(edge *GLPGEdge) {
	if edge == nil {
		return // Or handle error
	}
	g.Edges[edge.ID] = edge

	// Ensure source and target nodes exist in adjacency lists before appending
	if _, exists := g.OutgoingEdges[edge.SourceID]; !exists {
		// This case implies a node is being referenced before being added.
		// Depending on strictness, this could be an error or an implicit node creation.
		// For now, we'll assume nodes are added first or handle it gracefully.
		g.OutgoingEdges[edge.SourceID] = []*GLPGEdge{}
	}
	g.OutgoingEdges[edge.SourceID] = append(g.OutgoingEdges[edge.SourceID], edge)

	if _, exists := g.IncomingEdges[edge.TargetID]; !exists {
		g.IncomingEdges[edge.TargetID] = []*GLPGEdge{}
	}
	g.IncomingEdges[edge.TargetID] = append(g.IncomingEdges[edge.TargetID], edge)
}

// GetNode retrieves a node by its ID.
func (g *GLPG) GetNode(id string) *GLPGNode {
	return g.Nodes[id]
}

// GetEdge retrieves an edge by its ID.
func (g *GLPG) GetEdge(id string) *GLPGEdge {
	return g.Edges[id]
}

// GetOutgoingEdges retrieves all outgoing edges for a given node ID.
func (g *GLPG) GetOutgoingEdges(nodeID string) []*GLPGEdge {
	return g.OutgoingEdges[nodeID]
}

// GetIncomingEdges retrieves all incoming edges for a given node ID.
func (g *GLPG) GetIncomingEdges(nodeID string) []*GLPGEdge {
	return g.IncomingEdges[nodeID]
}
