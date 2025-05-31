package output

import (
	"fmt"
	"lazybox/internal/glpg"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/glamour"
)

// PrintGLPGAsMarkdown renders the GLPG in Markdown format.
func PrintGLPGAsMarkdown(graph *glpg.GLPG, flags map[string]bool) error {
	if graph == nil {
		fmt.Println("No data to render as Markdown.")
		return nil
	}

	var md strings.Builder

	md.WriteString(fmt.Sprintf("# GLPG Overview\n\n"))
	md.WriteString(fmt.Sprintf("**Nodes:** %d | **Edges:** %d\n\n", len(graph.Nodes), len(graph.Edges)))

	// Sort node IDs for consistent output
	nodeIDs := make([]string, 0, len(graph.Nodes))
	for id := range graph.Nodes {
		nodeIDs = append(nodeIDs, id)
	}
	sort.Strings(nodeIDs)

	for _, nodeID := range nodeIDs {
		node := graph.Nodes[nodeID]
		md.WriteString(renderNodeAsMarkdown(node, graph, flags))
		md.WriteString("\n---\n\n") // Separator between nodes
	}

	// Optionally, list all edges if verbose and no nodes (or if specifically requested)
	if verbose, _ := flags["verbose"]; verbose && len(graph.Edges) > 0 && len(graph.Nodes) == 0 {
		md.WriteString("## All Edges\n\n")
		md.WriteString("| ID | Source | Target | Label | Properties |\n")
		md.WriteString("|---|---|---|---|---|\n")
		// Sort edge IDs for consistent output
		edgeIDs := make([]string, 0, len(graph.Edges))
		for id := range graph.Edges {
			edgeIDs = append(edgeIDs, id)
		}
		sort.Strings(edgeIDs)
		for _, edgeID := range edgeIDs {
			edge := graph.Edges[edgeID]
			md.WriteString(renderEdgeAsMarkdownTableRow(edge))
		}
		md.WriteString("\n")
	}

	// Use glamour with a theme-aware style
	// Determine style based on theme (e.g., "dark", "light")
	// For now, using AutoStyle which tries to detect terminal theme.
	// We can enhance this to use specific glamour styles based on Base16.
	// For example, theme.GetDefaultTheme().GlamourStyle could be a string like "dark" or "light".
	// This requires mapping Base16 to glamour's style names or using custom glamour styles.
	// As a starting point, AutoStyle is reasonable.
	glamourStyle := glamour.WithAutoStyle()
	// If we had a mapping in our theme:
	// currentTheme := theme.GetDefaultTheme()
	// if currentTheme.GlamourStyle != "" { glamourStyle = glamour.WithStylePath(currentTheme.GlamourStyle) }

	renderer, err := glamour.NewTermRenderer(
		glamourStyle,
		glamour.WithWordWrap(100), // Adjust word wrap as needed
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to initialize glamour renderer:", err)
		fmt.Println(md.String()) // Print raw markdown if renderer fails
		return nil
	}

	out, err := renderer.Render(md.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to render markdown:", err)
		fmt.Println(md.String()) // Print raw markdown if rendering fails
		return nil
	}
	fmt.Print(out)
	return nil
}

func renderNodeAsMarkdown(node *glpg.GLPGNode, graph *glpg.GLPG, flags map[string]bool) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## Node: `%s`\n", node.ID))
	if len(node.Labels) > 0 {
		sb.WriteString(fmt.Sprintf("**Labels:** ` %s `\n", strings.Join(node.Labels, "`, `")))
	}
	sb.WriteString("\n")

	sb.WriteString("### Properties\n")
	if len(node.Properties) == 0 {
		sb.WriteString("_No properties_\n")
	} else {
		sb.WriteString("| Key | Value |\n")
		sb.WriteString("|---|---|\n")
		// Sort property keys for consistent output
		propKeys := make([]string, 0, len(node.Properties))
		for k := range node.Properties {
			propKeys = append(propKeys, k)
		}
		sort.Strings(propKeys)
		for _, key := range propKeys {
			val := node.Properties[key]
			// For string values that might contain markdown or code, wrap them in code blocks
			if strVal, ok := val.(string); ok && (strings.Contains(strVal, "\n") || strings.Contains(strVal, "`")) {
				sb.WriteString(fmt.Sprintf("| `%s` | ```\n%s\n``` |\n", key, strVal))
			} else {
				sb.WriteString(fmt.Sprintf("| `%s` | `%v` |\n", key, val))
			}
		}
	}
	sb.WriteString("\n")

	outgoing := graph.GetOutgoingEdges(node.ID)
	if len(outgoing) > 0 {
		sb.WriteString("### Outgoing Edges\n")
		sb.WriteString("| Label | Target Node | Properties |\n")
		sb.WriteString("|---|---|---|\n")
		for _, edge := range outgoing {
			propsStr := renderPropertiesAsMarkdownInline(edge.Properties)
			sb.WriteString(fmt.Sprintf("| `%s` | [`%s`](#node-%s) | %s |\n", edge.Label, edge.TargetID, edge.TargetID, propsStr))
		}
		sb.WriteString("\n")
	}

	incoming := graph.GetIncomingEdges(node.ID)
	if len(incoming) > 0 {
		sb.WriteString("### Incoming Edges\n")
		sb.WriteString("| Source Node | Label | Properties |\n")
		sb.WriteString("|---|---|---|\n")
		for _, edge := range incoming {
			propsStr := renderPropertiesAsMarkdownInline(edge.Properties)
			sb.WriteString(fmt.Sprintf("| [`%s`](#node-%s) | `%s` | %s |\n", edge.SourceID, edge.SourceID, edge.Label, propsStr))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func renderEdgeAsMarkdownTableRow(edge *glpg.GLPGEdge) string {
	propsStr := renderPropertiesAsMarkdownInline(edge.Properties)
	return fmt.Sprintf("| `%s` | [`%s`](#node-%s) | [`%s`](#node-%s) | `%s` | %s |\n",
		edge.ID, edge.SourceID, edge.SourceID, edge.TargetID, edge.TargetID, edge.Label, propsStr)
}

func renderPropertiesAsMarkdownInline(props glpg.GLPGProperty) string {
	if len(props) == 0 {
		return "_none_"
	}
	var parts []string
	// Sort property keys for consistent output
	propKeys := make([]string, 0, len(props))
	for k := range props {
		propKeys = append(propKeys, k)
	}
	sort.Strings(propKeys)

	for _, k := range propKeys {
		parts = append(parts, fmt.Sprintf("`%s`: `%v`", k, props[k]))
	}
	return strings.Join(parts, ", ")
}

/*
// Old PrintMarkdown function (to be removed or adapted)
func PrintMarkdown(info *ir.FileInfo) {
	if info == nil {
		fmt.Println("No data.")
		return
	}
	md := "# File Info\\n"
	md += fmt.Sprintf("- **Name:** %s\\n", info.Name)
	md += fmt.Sprintf("- **Path:** %s\\n", info.Path)
	md += fmt.Sprintf("- **Type:** %s\\n", info.Type)
	md += fmt.Sprintf("- **Size:** %d bytes\\n", info.Size)
	if info.Type == "text" {
		if info.LineCount > 0 {
			md += fmt.Sprintf("- **Line Count:** %d\\n", info.LineCount)
		}
		if info.WordCount > 0 {
			md += fmt.Sprintf("- **Word Count:** %d\\n", info.WordCount)
		}
	}
	if info.Content != "" {
		md += "\\n## Content\\n\\n```" + string(info.Type) + "\\n" + info.Content + "\\n```\\n"
	}
	if len(info.Contents) > 0 {
		md += "\\n## Directory Contents\\n\\n| Name | Type | Size |\\n|---|---|---|\\n"
		for _, c := range info.Contents {
			md += fmt.Sprintf("| %s | %s | %d |\\n", c.Name, c.Type, c.Size)
		}
	}
	// Use glamour with a high-contrast theme and syntax highlighting
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to initialize glamour renderer:", err)
		fmt.Println(md)
		return
	}
	out, err := renderer.Render(md)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to render markdown:", err)
		fmt.Println(md)
		return
	}
	fmt.Print(out)
}
*/
