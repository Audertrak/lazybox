package output

import (
	"fmt"
	"os"
	"strings"

	"lazybox/internal/ir"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(themeColor(CurrentTheme.Base0D)).Background(themeColor(CurrentTheme.Base01)).Padding(0, 1)
	labelStyle  = lipgloss.NewStyle().Bold(true).Foreground(themeColor(CurrentTheme.Base0B))
	valueStyle  = lipgloss.NewStyle().Foreground(themeColor(CurrentTheme.Base05))
	boxStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(themeColor(CurrentTheme.Base0C)).Padding(1, 2).Margin(1, 0).Background(themeColor(CurrentTheme.Base00))
	codeStyle   = lipgloss.NewStyle().Foreground(themeColor(CurrentTheme.Base0A)).Background(themeColor(CurrentTheme.Base01)).Padding(0, 1)
)

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
		// Use glamour for code block rendering
		md := "```" + string(info.Type) + "\n" + info.Content + "\n```"
		out, err := glamour.Render(md, "notty")
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
	os.Stdout.Sync()
}
