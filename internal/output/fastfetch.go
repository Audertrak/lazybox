package output

import (
	"fmt"
	"lazybox/internal/glpg"
	"lazybox/internal/theme"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// PrintGLPGAsFastfetch renders a GLPG in a style inspired by fastfetch/neofetch.
// It displays a summary of the graph, including node and edge counts, and then
// lists nodes and their properties, styled with the current Base16 theme.
func PrintGLPGAsFastfetch(graph *glpg.GLPG, flags map[string]bool) error {
	if graph == nil {
		// TODO: Themed output for "no data"
		fmt.Println("No data to display.")
		return nil
	}

	currentTheme := theme.GetDefaultTheme() // Assuming theme is initialized

	// Define styles based on the theme
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(currentTheme.Base0D))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(currentTheme.Base0B)).Bold(true)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(currentTheme.Base05))
	separatorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(currentTheme.Base03))
	containerStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(currentTheme.Base0C)).Padding(1, 2).MarginBottom(1)

	var sb strings.Builder

	// --- System Info / Graph Summary ---
	osName := runtime.GOOS
	arch := runtime.GOARCH
	user := os.Getenv("USER") // Simpler way to get user, fallback if empty
	if user == "" {
		user = os.Getenv("USERNAME")
	}
	host, _ := os.Hostname()
	goVersion := runtime.Version()

	sb.WriteString(titleStyle.Render("lazybox GLPG Report") + "\n")
	sb.WriteString(labelStyle.Render("User: ") + valueStyle.Render(fmt.Sprintf("%s@%s", user, host)) + "\n")
	sb.WriteString(labelStyle.Render("OS:   ") + valueStyle.Render(osName) + "\n")
	sb.WriteString(labelStyle.Render("Arch: ") + valueStyle.Render(arch) + "\n")
	sb.WriteString(labelStyle.Render("Go:   ") + valueStyle.Render(goVersion) + "\n")
	sb.WriteString(labelStyle.Render("Nodes: ") + valueStyle.Render(fmt.Sprintf("%d", len(graph.Nodes))) + "\n")
	sb.WriteString(labelStyle.Render("Edges: ") + valueStyle.Render(fmt.Sprintf("%d", len(graph.Edges))) + "\n")
	sb.WriteString(separatorStyle.Render(strings.Repeat("─", 40)) + "\n\n")

	// --- Nodes Details ---
	sb.WriteString(titleStyle.Render("Nodes") + "\n")

	// Sort node IDs for consistent output
	nodeIDs := make([]string, 0, len(graph.Nodes))
	for id := range graph.Nodes {
		nodeIDs = append(nodeIDs, id)
	}
	sort.Strings(nodeIDs)

	for _, nodeID := range nodeIDs {
		node := graph.Nodes[nodeID]
		nodeLabelStr := ""
		if len(node.Labels) > 0 {
			nodeLabelStr = strings.Join(node.Labels, ", ") // Join labels if multiple, or use the first one
		}
		sb.WriteString(labelStyle.Render("ID:    ") + valueStyle.Render(node.ID) + "\n")
		sb.WriteString(labelStyle.Render("Labels: ") + valueStyle.Render(nodeLabelStr) + "\n") // Changed from Label to Labels

		// Sort property keys for consistent output
		propKeys := make([]string, 0, len(node.Properties))
		for k := range node.Properties {
			propKeys = append(propKeys, k)
		}
		sort.Strings(propKeys)

		for _, key := range propKeys {
			val := node.Properties[key]
			sb.WriteString(labelStyle.Render(fmt.Sprintf("  %s: ", key)) + valueStyle.Render(fmt.Sprintf("%v", val)) + "\n")
		}
		sb.WriteString("\n")
	}

	// --- Edges Details (Optional, can be verbose) ---
	if flags["edges"] { // Example flag to control edge printing
		sb.WriteString(separatorStyle.Render(strings.Repeat("─", 40)) + "\n\n")
		sb.WriteString(titleStyle.Render("Edges") + "\n")

		// Sort edge IDs for consistent output
		edgeIDs := make([]string, 0, len(graph.Edges))
		for id := range graph.Edges {
			edgeIDs = append(edgeIDs, id)
		}
		sort.Strings(edgeIDs)

		for _, edgeID := range edgeIDs {
			edge := graph.Edges[edgeID]
			sb.WriteString(labelStyle.Render("ID:     ") + valueStyle.Render(edge.ID) + "\n")
			sb.WriteString(labelStyle.Render("Label:  ") + valueStyle.Render(edge.Label) + "\n")
			sb.WriteString(labelStyle.Render("Source: ") + valueStyle.Render(edge.SourceID) + "\n")
			sb.WriteString(labelStyle.Render("Target: ") + valueStyle.Render(edge.TargetID) + "\n")

			// Sort property keys for consistent output
			propKeys := make([]string, 0, len(edge.Properties))
			for k := range edge.Properties {
				propKeys = append(propKeys, k)
			}
			sort.Strings(propKeys)

			for _, key := range propKeys {
				val := edge.Properties[key]
				sb.WriteString(labelStyle.Render(fmt.Sprintf("  %s: ", key)) + valueStyle.Render(fmt.Sprintf("%v", val)) + "\n")
			}
			sb.WriteString("\n")
		}
	}

	fmt.Println(containerStyle.Render(sb.String()))
	return nil
}

// themeColor is a helper, assuming it's defined elsewhere or we use lipgloss.Color directly
// For this refactor, we'll use lipgloss.Color(currentTheme.BaseXX) directly.

// getEnv and related functions are removed as direct os.Getenv or os.Hostname is simpler for this context.
// The original fastfetch had a more complex env var retrieval, which is not strictly needed for GLPG display.
