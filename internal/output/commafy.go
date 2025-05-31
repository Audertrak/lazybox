package output

import (
	"fmt"
	"os"
	"strings"

	"lazybox/internal/glpg"
	"lazybox/internal/theme"

	"github.com/charmbracelet/lipgloss"
)

// PrintGLPGAsCommafy prints the GLPG data as a comma-separated list of property values.
func PrintGLPGAsCommafy(g *glpg.GLPG, flags map[string]bool) error {
	if g == nil {
		return fmt.Errorf("cannot print nil GLPG")
	}

	th := theme.GetDefaultTheme()
	var values []string

	style := lipgloss.NewStyle().Foreground(lipgloss.Color(string(th.Base05))) // Corrected: Use th.Base05

	// Collect node property values
	for _, node := range g.Nodes {
		for key, value := range node.Properties { // Iterate over map
			values = append(values, fmt.Sprintf("%v:%v", key, value))
		}
	}

	// Collect edge property values
	for _, edge := range g.Edges {
		for key, value := range edge.Properties { // Iterate over map
			values = append(values, fmt.Sprintf("%v:%v", key, value))
		}
	}

	if len(values) == 0 {
		fmt.Fprintln(os.Stdout, style.Render("No properties to display in comma-separated format."))
		return nil
	}

	outputString := strings.Join(values, ", ")
	fmt.Fprintln(os.Stdout, style.Render(outputString))

	return nil
}
