package output

import (
	"fmt"
	"lazybox/internal/ir"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

var (
	tableHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(themeColor(CurrentTheme.Base0D))
	rowStyleEven     = lipgloss.NewStyle().Foreground(themeColor(CurrentTheme.Base05))
	rowStyleOdd      = lipgloss.NewStyle().Foreground(themeColor(CurrentTheme.Base03))
)

func PrintTable(info *ir.FileInfo) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.Style().Color.Header = text.Colors{text.FgCyan, text.Bold}
	t.Style().Color.Row = text.Colors{text.FgWhite}
	t.Style().Color.RowAlternate = text.Colors{text.FgHiBlack}
	t.Style().Options.DrawBorder = true
	t.Style().Options.SeparateRows = true

	if info == nil {
		fmt.Println(tableHeaderStyle.Render("No data."))
		return
	}

	if info.Type == ir.FileTypeDirectory {
		t.AppendHeader(table.Row{"Name", "Type", "Size", "Extension"})
		for _, c := range info.Contents {
			row := table.Row{c.Name, string(c.Type), c.Size, c.Extension}
			t.AppendRow(row)
		}
		t.Render()
		return
	}

	t.AppendHeader(table.Row{"Field", "Value"})
	metaRows := [][]string{
		{"Name", info.Name},
		{"Path", info.Path},
		{"Type", string(info.Type)},
		{"Size", fmt.Sprint(info.Size)},
		{"Extension", info.Extension},
	}
	if info.LineCount > 0 {
		metaRows = append(metaRows, []string{"LineCount", fmt.Sprint(info.LineCount)})
	}
	if info.WordCount > 0 {
		metaRows = append(metaRows, []string{"WordCount", fmt.Sprint(info.WordCount)})
	}
	for _, row := range metaRows {
		t.AppendRow(table.Row{row[0], row[1]})
	}
	t.Render()
	if info.Content != "" {
		fmt.Println(tableHeaderStyle.Render("\nContent Preview:\n----------------"))
		fmt.Println(rowStyleEven.Render(info.Content))
	}
}
