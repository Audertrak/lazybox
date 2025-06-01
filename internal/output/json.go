package output

import (
	"encoding/json"
	"fmt"
	"lazybox/internal/glpg"
	"lazybox/internal/theme"
	"os"
	"regexp"

	"github.com/charmbracelet/lipgloss"
)

var (
	fillerWords = regexp.MustCompile(`\b(a|an|the|very|just|basically|especially|really|actually|literally|simply|uh+|um+|like|so|well|sort of|kind of|maybe|perhaps|that is|in order to|in other words|for example|for instance|such as|etc|etc\.|and so on|and so forth)\b`)
	phraseSubs  = map[string]string{
		"approximately":        "~",
		"information":          "info",
		"configuration":        "config",
		"application":          "app",
		"directory":            "dir",
		"function":             "fn",
		"parameter":            "param",
		"argument":             "arg",
		"representation":       "repr",
		"documentation":        "docs",
		"command line":         "CLI",
		"large language model": "LLM",
		"language model":       "LM",
		"with respect to":      "re: ",
		"as soon as possible":  "ASAP",
		"for your information": "FYI",
		"by the way":           "BTW",
		"see also":             "cf.",
		"for example":          "e.g.",
		"that is":              "i.e.",
		"and so on":            "...",
		"and so forth":         "...",
	}
	jqKeyStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.GetDefaultTheme().Base0B)).Bold(true)
	jqStringStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.GetDefaultTheme().Base0A))
	jqNumStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.GetDefaultTheme().Base09))
	jqBoolStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.GetDefaultTheme().Base0E)).Bold(true)
	jqNullStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.GetDefaultTheme().Base03)).Italic(true)
)

// PrintGLPGAsJSON serializes the GLPG to a styled JSON output (always styled, jq-inspired)
func PrintGLPGAsJSON(graph *glpg.GLPG, flags map[string]bool) error {
	if minFlag, _ := flags["min"]; minFlag {
		return PrintGLPGAsMinJSON(graph)
	}
	if lessFlag, _ := flags["less"]; lessFlag {
		return PrintGLPGAsLessJSON(graph)
	}
	data, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling GLPG to JSON: %v\n", err)
		return err
	}
	// Always style the output (jq-inspired coloring)
	coloredOutput := string(data)
	coloredOutput = regexp.MustCompile(`"([^"]+)":`).ReplaceAllStringFunc(coloredOutput, func(match string) string {
		parts := regexp.MustCompile(`"([^"]+)":`).FindStringSubmatch(match)
		if len(parts) > 1 {
			return jqKeyStyle.Render(fmt.Sprintf(`"%s"`, parts[1])) + ":"
		}
		return match
	})
	coloredOutput = regexp.MustCompile(`: "([^"]*)"`).ReplaceAllStringFunc(coloredOutput, func(match string) string {
		parts := regexp.MustCompile(`: "([^"]*)"`).FindStringSubmatch(match)
		if len(parts) > 1 {
			return ": " + jqStringStyle.Render(fmt.Sprintf(`"%s"`, parts[1]))
		}
		return match
	})
	coloredOutput = regexp.MustCompile(`: (\d+\.?\d*)`).ReplaceAllStringFunc(coloredOutput, func(match string) string {
		parts := regexp.MustCompile(`: (\d+\.?\d*)`).FindStringSubmatch(match)
		if len(parts) > 1 {
			return ": " + jqNumStyle.Render(parts[1])
		}
		return match
	})
	coloredOutput = regexp.MustCompile(`: (true|false)`).ReplaceAllStringFunc(coloredOutput, func(match string) string {
		parts := regexp.MustCompile(`: (true|false)`).FindStringSubmatch(match)
		if len(parts) > 1 {
			return ": " + jqBoolStyle.Render(parts[1])
		}
		return match
	})
	coloredOutput = regexp.MustCompile(`: null`).ReplaceAllStringFunc(coloredOutput, func(match string) string {
		return ": " + jqNullStyle.Render("null")
	})
	fmt.Println(coloredOutput)
	return nil
}

// PrintGLPGAsMinJSON serializes the GLPG to a minified JSON output.
func PrintGLPGAsMinJSON(graph *glpg.GLPG) error {
	data, err := json.Marshal(graph)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling GLPG to minified JSON: %v\\n", err)
		return err
	}
	fmt.Println(string(data))
	return nil
}

// PrintGLPGAsLessJSON provides a summarized JSON view of the GLPG.
// This might involve showing counts of nodes/edges or specific high-level properties.
func PrintGLPGAsLessJSON(graph *glpg.GLPG) error {
	if graph == nil {
		fmt.Println("{}")
		return nil
	}
	summary := struct {
		NodeCount int      `json:"node_count"`
		EdgeCount int      `json:"edge_count"`
		NodeIDs   []string `json:"node_ids,omitempty"` // Example: show some node IDs
		// Potentially add more summary fields here, like label counts, etc.
	}{
		NodeCount: len(graph.Nodes),
		EdgeCount: len(graph.Edges),
	}

	// Optionally, add a few node IDs to the summary
	maxPreviewIDs := 5
	for id := range graph.Nodes {
		if len(summary.NodeIDs) < maxPreviewIDs {
			summary.NodeIDs = append(summary.NodeIDs, id)
		} else {
			break
		}
	}

	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling GLPG summary to JSON: %v\\n", err)
		return err
	}
	fmt.Println(string(data))
	return nil
}
