package output

import (
	"fmt"
	"lazybox/internal/glpg"
	"lazybox/internal/theme"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

// PrintGLPGAsFlow renders the GLPG as a styled flowchart/diagram (ASCII art, theme-styled)
func PrintGLPGAsFlow(data *glpg.GLPG, flags map[string]bool) error {
	ct := theme.GetDefaultTheme()
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ct.Base0D)).
		Foreground(lipgloss.Color(ct.Base05)).
		Padding(0, 1)
	arrow := lipgloss.NewStyle().Foreground(lipgloss.Color(ct.Base0E)).Render("→")

title := lipgloss.NewStyle().Foreground(lipgloss.Color(ct.Base0A)).Bold(true).Render("GLPG Flowchart")

	var b strings.Builder
	b.WriteString(title + "\n\n")

	nodeIDs := make([]string, 0, len(data.Nodes))
	for id := range data.Nodes {
		nodeIDs = append(nodeIDs, id)
	}

	// For each node, print the box and its outgoing edges as arrows to other boxes
	for _, nodeID := range nodeIDs {
		node := data.Nodes[nodeID]
		b.WriteString(box.Render(node.ID))
		outEdges := data.GetOutgoingEdges(nodeID)
		for _, edge := range outEdges {
			target := data.GetNode(edge.TargetID)
			if target != nil {
				b.WriteString("\n  " + arrow + " " + box.Render(target.ID) + "\n")
			}
		}
		b.WriteString("\n")
	}

	fmt.Println(b.String())
	return nil
}
