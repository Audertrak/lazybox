package output

import (
	"fmt"
	"lazybox/internal/glpg"
	"lazybox/internal/theme"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// Styles (can be initialized in an init() function or when first used)
var (
	headerStyle      lipgloss.Style
	labelStyle       lipgloss.Style
	valueStyle       lipgloss.Style
	boxStyle         lipgloss.Style
	codeStyle        lipgloss.Style
	nodeIDStyle      lipgloss.Style
	edgeLabelStyle   lipgloss.Style
	propertyKeyStyle lipgloss.Style
	tableHeaderStyle lipgloss.Style
	tableCellStyle   lipgloss.Style
)

func initializeStyles() {
	ct := theme.GetDefaultTheme()
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ct.Base0D)).Background(lipgloss.Color(ct.Base01)).Padding(0, 1).MarginBottom(1)
	labelStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ct.Base0B)).MarginRight(1)
	valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.Base05))
	boxStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(ct.Base0C)).Padding(1).Margin(0, 1, 1, 1)
	codeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.Base0A)).Background(lipgloss.Color(ct.Base01)).Padding(0, 1)
	nodeIDStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.Base0E)).Bold(true)
	edgeLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.Base0C)).Italic(true)
	propertyKeyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.Base03))
	// tableHeaderStyle is defined globally
	tableHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ct.Base05)).BorderBottom(true).BorderForeground(lipgloss.Color(ct.Base03))
	tableCellStyle = lipgloss.NewStyle().Padding(0, 1)
}

func init() {
	initializeStyles()
}

// PrintGLPGAsPretty renders the GLPG in a human-readable, styled format.
func PrintGLPGAsPretty(graph *glpg.GLPG, flags map[string]bool) error {
	if graph == nil {
		fmt.Println(boxStyle.Render(headerStyle.Render("Empty Graph")))
		return nil
	}

	var b strings.Builder

	b.WriteString(headerStyle.Render(fmt.Sprintf("GLPG Overview (Nodes: %d, Edges: %d)", len(graph.Nodes), len(graph.Edges))))
	b.WriteString("\n")

	// Sort node IDs for consistent output
	nodeIDs := make([]string, 0, len(graph.Nodes))
	for id := range graph.Nodes {
		nodeIDs = append(nodeIDs, id)
	}
	sort.Strings(nodeIDs)

	for _, nodeID := range nodeIDs {
		node := graph.Nodes[nodeID]
		nodeBox := renderNode(node, graph, flags)
		b.WriteString(boxStyle.Render(nodeBox))
		b.WriteString("\n")
	}

	// For verbose, could also list all edges separately if not covered by node outgoing/incoming
	if verbose, _ := flags["verbose"]; verbose && len(graph.Edges) > 0 && len(graph.Nodes) == 0 {
		b.WriteString(headerStyle.Render("All Edges"))
		b.WriteString("\n")
		edgesTable := createEdgesTable(graph.Edges, "All Edges Overview") // Title passed for context, not set on table directly
		// Render the table using PrintTableAsString
		// Note: PrintTableAsString currently prints directly.
		// For building a single string `b`, we might need createEdgesTable to return a string
		// or PrintTableAsString to return a string.
		// For now, let's assume we want to print it directly if this condition is met.
		// To integrate into the string builder `b`, PrintTableAsString would need to change.
		// As a temporary measure, we'll print it directly here.
		// A better approach would be for PrintTableAsString to return the string.
		fmt.Println(b.String())        // Print what we have so far
		PrintTableAsString(edgesTable) // Print the table
		b.Reset()                      // Reset builder as table was printed separately
	} else {
		fmt.Println(b.String())
	}
	return nil
}

func renderNode(node *glpg.GLPGNode, graph *glpg.GLPG, flags map[string]bool) string {
	var sb strings.Builder
	sb.WriteString(nodeIDStyle.Render(fmt.Sprintf("Node: %s", node.ID)))
	if len(node.Labels) > 0 {
		sb.WriteString(fmt.Sprintf(" (%s)", valueStyle.Render(strings.Join(node.Labels, ", "))))
	}
	sb.WriteString("\n")

	sb.WriteString(labelStyle.Render("Properties:") + "\n")
	if len(node.Properties) == 0 {
		sb.WriteString(valueStyle.Italic(true).Render("  <none>") + "\n")
	} else {
		// Sort property keys for consistent output
		propKeys := make([]string, 0, len(node.Properties))
		for k := range node.Properties {
			propKeys = append(propKeys, k)
		}
		sort.Strings(propKeys)
		for _, key := range propKeys {
			val := node.Properties[key]
			// Attempt to render content with glamour if it looks like markdown or code
			// This is a heuristic
			if strVal, ok := val.(string); ok && (strings.Contains(strVal, "\n") || strings.HasPrefix(strVal, "```")) {
				mdVal, err := glamour.Render(fmt.Sprintf("```\n%s\n```", strVal), "dark") // or use theme
				if err == nil {
					sb.WriteString(fmt.Sprintf("  %s %v\n", propertyKeyStyle.Render(key+":"), mdVal))
				} else {
					sb.WriteString(fmt.Sprintf("  %s %s\n", propertyKeyStyle.Render(key+":"), codeStyle.Render(strVal)))
				}
			} else {
				sb.WriteString(fmt.Sprintf("  %s %s\n", propertyKeyStyle.Render(key+":"), valueStyle.Render(fmt.Sprintf("%v", val))))
			}
		}
	}

	outgoing := graph.GetOutgoingEdges(node.ID)
	if len(outgoing) > 0 {
		sb.WriteString(labelStyle.Render("Outgoing Edges:") + "\n")
		for _, edge := range outgoing {
			sb.WriteString(fmt.Sprintf("  %s %s %s\n",
				edgeLabelStyle.Render("--("+edge.Label+")-->"),
				nodeIDStyle.Render(edge.TargetID),
				renderPropertiesInline(edge.Properties),
			))
		}
	}

	incoming := graph.GetIncomingEdges(node.ID)
	if len(incoming) > 0 {
		sb.WriteString(labelStyle.Render("Incoming Edges:") + "\n")
		for _, edge := range incoming {
			sb.WriteString(fmt.Sprintf("  %s %s %s\n",
				nodeIDStyle.Render(edge.SourceID),
				edgeLabelStyle.Render("--("+edge.Label+")-->"),
				renderPropertiesInline(edge.Properties),
			))
		}
	}

	return sb.String()
}

func renderEdgeSimple(edge *glpg.GLPGEdge) string {
	return fmt.Sprintf("%s %s %s %s %s",
		nodeIDStyle.Render(edge.SourceID),
		edgeLabelStyle.Render("--("+edge.Label+")-->"),
		nodeIDStyle.Render(edge.TargetID),
		valueStyle.Render("ID: "+edge.ID),
		renderPropertiesInline(edge.Properties),
	)
}

func renderPropertiesInline(props glpg.GLPGProperty) string {
	if len(props) == 0 {
		return ""
	}
	var parts []string
	// Sort property keys for consistent output
	propKeys := make([]string, 0, len(props))
	for k := range props {
		propKeys = append(propKeys, k)
	}
	sort.Strings(propKeys)

	for _, k := range propKeys {
		parts = append(parts, fmt.Sprintf("%s:%s", propertyKeyStyle.Render(k), valueStyle.Render(fmt.Sprintf("%v", props[k]))))
	}
	return "(" + strings.Join(parts, ", ") + ")"
}

// createEdgesTable prepares a bubbletea table model for edges.
// This is more for interactive TUIs. For static output, a simpler format might be better.
func createEdgesTable(edges map[string]*glpg.GLPGEdge, title string) table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 15},
		{Title: "Source", Width: 15},
		{Title: "Target", Width: 15},
		{Title: "Label", Width: 15},
		{Title: "Properties", Width: 30},
	}

	var rows []table.Row
	// Sort edge IDs for consistent output
	edgeIDs := make([]string, 0, len(edges))
	for id := range edges {
		edgeIDs = append(edgeIDs, id)
	}
	sort.Strings(edgeIDs)

	for _, id := range edgeIDs {
		edge := edges[id]
		propsStr := renderPropertiesInline(edge.Properties)
		rows = append(rows, table.Row{edge.ID, edge.SourceID, edge.TargetID, edge.Label, propsStr})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	// ct := theme.GetDefaultTheme() // ct is already available in initializeStyles, or pass theme explicitly
	s.Header = tableHeaderStyle // Use the globally defined and initialized style
	s.Cell = tableCellStyle
	// s.Selected = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color(ct.Base02)).Foreground(lipgloss.Color(ct.Base05)) // ct needs to be accessible here
	// For s.Selected, ensure ct is accessible or pass it. For now, let's re-fetch or pass.
	currentTheme := theme.GetDefaultTheme()
	s.Selected = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color(currentTheme.Base02)).Foreground(lipgloss.Color(currentTheme.Base05))
	t.SetStyles(s)
	// t.Title = title // table.Model does not have a Title field. Title should be rendered separately.

	return t
}

// Since bubbletea is for interactive apps, we need a way to print the table model to string for static output.
// This is a simplified approach. Bubbletea tables are best used in a tea.Program.
func PrintTableAsString(m table.Model) {
	// For a static CLI, we just print the view.
	// A real bubbletea app would run p := tea.NewProgram(m); p.Start()
	// If a title was intended for the table, print it here above the table.
	// This function currently doesn't know the title.
	// We could modify createEdgesTable to also return the title string,
	// or pass the title to this function.
	// For now, the title is rendered as a headerStyle string in PrintGLPGAsPretty.
	fmt.Println(m.View())
}

// Old PrintPretty function (to be removed or adapted if specific IR.FileInfo formatting is still needed elsewhere)
/*
	func PrintPretty(info *ir.FileInfo) {
		if info == nil {
			fmt.Println("No data.")
			return
		}
		var b strings.Builder
		b.WriteString(headerStyle.Render(" File Info "))
		b.WriteString("\n")
		b.WriteString(labelStyle.Render("Name: ") + valueStyle.Render(info.Name) + "\n")
		b.WriteString(labelStyle.Render("Path: ") + valueStyle.Render(info.Path) + "\n")
		b.WriteString(labelStyle.Render("Type: ") + valueStyle.Render(string(info.Type)) + "\n")
		b.WriteString(labelStyle.Render("Size: ") + valueStyle.Render(fmt.Sprintf("%d bytes", info.Size)) + "\n")
		if info.Type == "text" {
			if info.LineCount > 0 {
				b.WriteString(labelStyle.Render("Line Count: ") + valueStyle.Render(fmt.Sprintf("%d", info.LineCount)) + "\n")
			}
			if info.WordCount > 0 {
				b.WriteString(labelStyle.Render("Word Count: ") + valueStyle.Render(fmt.Sprintf("%d", info.WordCount)) + "\n")
			}
		}
		if info.Content != "" {
			b.WriteString("\n" + labelStyle.Render("Content:") + "\n")
			md := "```" + string(info.Type) + "\n" + info.Content + "\n```"
			out, err := glamour.Render(md, "notty") // Use "notty" for non-interactive rendering
			if err == nil {
				b.WriteString(out)
			} else {
				b.WriteString(codeStyle.Render(info.Content) + "\n")
			}
		}
		if len(info.Contents) > 0 {
			b.WriteString("\n" + labelStyle.Render("Directory Contents:") + "\n")
			for _, c := range info.Contents {
				row := fmt.Sprintf("%s  %s  %d\n", valueStyle.Render(c.Name), valueStyle.Render(string(c.Type)), c.Size)
				b.WriteString(row)
			}
		}
		fmt.Print(boxStyle.Render(b.String()))
		os.Stdout.Sync() // Ensure output is flushed
	}
*/
