package output

import (
	"encoding/json"
	"fmt"
	"lazybox/internal/glpg"
	"lazybox/internal/theme"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Supported comment styles for various languages
var commentDelimiters = map[string][2]string{
	"bash":   {"# ", ""},
	"python": {"# ", ""},
	"go":     {"// ", ""},
	"c":      {"/* ", " */"},
	"lua":    {"-- ", ""},
	"sql":    {"-- ", ""},
}

// PrintGLPGAsComment renders the GLPG as a code comment block in the specified language
func PrintGLPGAsComment(data *glpg.GLPG, flags map[string]bool, lang string) error {
	ct := theme.GetDefaultTheme()
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(ct.Base05))
	border := lipgloss.NewStyle().Foreground(lipgloss.Color(ct.Base0E)).Bold(true)

	// Pick comment delimiters
	delims, ok := commentDelimiters[lang]
	if !ok {
		delims = commentDelimiters["bash"] // fallback
	}

	// Marshal GLPG to pretty JSON for comment block
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	jsonLines := strings.Split(string(jsonData), "\n")

	var b strings.Builder
	b.WriteString(border.Render("Commentified Output ("+lang+")") + "\n")
	b.WriteString("\n")
	for _, line := range jsonLines {
		b.WriteString(style.Render(delims[0]+line+delims[1]) + "\n")
	}

	fmt.Println(b.String())
	return nil
}
