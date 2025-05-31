package output

import (
	"fmt"
	"lazybox/internal/glpg"
	"lazybox/internal/theme"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// Styles for table output - initialized in init()
var (
	tblHeaderStyle       lipgloss.Style
	tblCellStyle         lipgloss.Style
	tblSelectedCellStyle lipgloss.Style
	tblContainerStyle    lipgloss.Style
	tblTitleStyle        lipgloss.Style
)

func initializeTableStyles() {
	ct := theme.GetDefaultTheme()
	tblHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ct.Base0D)).Padding(0, 1).BorderBottom(true).BorderForeground(lipgloss.Color(ct.Base03))
	tblCellStyle = lipgloss.NewStyle().Padding(0, 1)
	tblSelectedCellStyle = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color(ct.Base02)).Foreground(lipgloss.Color(ct.Base05))
	tblContainerStyle = lipgloss.NewStyle().Margin(1, 0)
	tblTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ct.Base0E)).MarginBottom(1)
}

func init() {
	initializeTableStyles() // Initialize styles when package is loaded
}

// PrintGLPGAsTable renders the GLPG data as styled tables using Charmbracelet components.
func PrintGLPGAsTable(graph *glpg.GLPG, flags map[string]bool) error {
	if graph == nil || (len(graph.Nodes) == 0 && len(graph.Edges) == 0) {
		fmt.Println(tblContainerStyle.Render(tblTitleStyle.Render("No data to display in table.")))
		return nil
	}

	var output strings.Builder

	// Render Nodes Table
	if len(graph.Nodes) > 0 {
		nodesTableStr := renderNodesTable(graph, flags)
		output.WriteString(nodesTableStr)
		output.WriteString("\n")
	}

	// Render Edges Table
	if len(graph.Edges) > 0 {
		edgesTableStr := renderEdgesTable(graph, flags)
		output.WriteString(edgesTableStr)
	}

	fmt.Print(output.String())
	return nil
}

func renderNodesTable(graph *glpg.GLPG, flags map[string]bool) string {
	columns := []table.Column{
		{Title: "Node ID", Width: 20},
		{Title: "Labels", Width: 25},
		{Title: "Properties", Width: 50},
	}

	var rows []table.Row
	nodeIDs := make([]string, 0, len(graph.Nodes))
	for id := range graph.Nodes {
		nodeIDs = append(nodeIDs, id)
	}
	sort.Strings(nodeIDs) // Sort for consistent output

	for _, id := range nodeIDs {
		node := graph.Nodes[id]
		labelsStr := strings.Join(node.Labels, ", ")
		propsStr := formatPropertiesForTable(node.Properties)
		rows = append(rows, table.Row{id, labelsStr, propsStr})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),      // No focus needed for static display
		table.WithHeight(len(rows)+1), // Adjust height dynamically
		table.WithStyles(table.Styles{
			Header:   tblHeaderStyle,
			Cell:     tblCellStyle,
			Selected: tblSelectedCellStyle,
		}),
	)

	return tblContainerStyle.Render(tblTitleStyle.Render("Nodes") + "\n" + t.View())
}

func renderEdgesTable(graph *glpg.GLPG, flags map[string]bool) string {
	columns := []table.Column{
		{Title: "Edge ID", Width: 20},
		{Title: "Source ID", Width: 20},
		{Title: "Target ID", Width: 20},
		{Title: "Label", Width: 20},
		{Title: "Properties", Width: 35},
	}

	var rows []table.Row
	edgeIDs := make([]string, 0, len(graph.Edges))
	for id := range graph.Edges {
		edgeIDs = append(edgeIDs, id)
	}
	sort.Strings(edgeIDs) // Sort for consistent output

	for _, id := range edgeIDs {
		edge := graph.Edges[id]
		propsStr := formatPropertiesForTable(edge.Properties)
		rows = append(rows, table.Row{id, edge.SourceID, edge.TargetID, edge.Label, propsStr})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(len(rows)+1),
		table.WithStyles(table.Styles{
			Header:   tblHeaderStyle,
			Cell:     tblCellStyle,
			Selected: tblSelectedCellStyle,
		}),
	)

	return tblContainerStyle.Render(tblTitleStyle.Render("Edges") + "\n" + t.View())
}

// formatPropertiesForTable converts GLPGProperty map to a string for table display.
func formatPropertiesForTable(props glpg.GLPGProperty) string {
	if len(props) == 0 {
		return "-"
	}
	var parts []string
	// Sort keys for consistent output
	keys := make([]string, 0, len(props))
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := props[k]
		// Simple string representation; could be truncated or ellipsized if too long
		valStr := fmt.Sprintf("%v", v)
		if len(valStr) > 30 { // Arbitrary limit for inline display
			valStr = valStr[:27] + "..."
		}
		parts = append(parts, fmt.Sprintf("%s: %s", k, valStr))
	}
	return strings.Join(parts, "; ")
}
