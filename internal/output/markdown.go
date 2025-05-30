package output

import (
	"fmt"
	"lazybox/internal/ir"
	"os"

	"github.com/charmbracelet/glamour"
)

func PrintMarkdown(info *ir.FileInfo) {
	if info == nil {
		fmt.Println("No data.")
		return
	}
	md := "# File Info\n"
	md += fmt.Sprintf("- **Name:** %s\n", info.Name)
	md += fmt.Sprintf("- **Path:** %s\n", info.Path)
	md += fmt.Sprintf("- **Type:** %s\n", info.Type)
	md += fmt.Sprintf("- **Size:** %d bytes\n", info.Size)
	if info.Type == "text" {
		if info.LineCount > 0 {
			md += fmt.Sprintf("- **Line Count:** %d\n", info.LineCount)
		}
		if info.WordCount > 0 {
			md += fmt.Sprintf("- **Word Count:** %d\n", info.WordCount)
		}
	}
	if info.Content != "" {
		md += "\n## Content\n\n```" + string(info.Type) + "\n" + info.Content + "\n```\n"
	}
	if len(info.Contents) > 0 {
		md += "\n## Directory Contents\n\n| Name | Type | Size |\n|---|---|---|\n"
		for _, c := range info.Contents {
			md += fmt.Sprintf("| %s | %s | %d |\n", c.Name, c.Type, c.Size)
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
