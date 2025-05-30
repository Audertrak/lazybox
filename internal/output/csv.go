package output

import (
	"encoding/csv"
	"fmt"
	"lazybox/internal/ir"
	"os"
)

func PrintCSV(info *ir.FileInfo) {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	if info == nil {
		fmt.Println("No data.")
		return
	}

	// Print header and row for file or text/code targets
	if info.Type == "file" || info.Type == "text" || info.Type == "code" || info.Type == "func" || info.Type == "struct" || info.Type == "enum" || info.Type == "list" {
		head := []string{"Name", "Type", "Size", "Extension", "LineCount", "WordCount"}
		row := []string{info.Name, string(info.Type), fmt.Sprint(info.Size), info.Extension, fmt.Sprint(info.LineCount), fmt.Sprint(info.WordCount)}
		w.Write(head)
		w.Write(row)
		if info.Content != "" {
			w.Write([]string{"Content", info.Content})
		}
		return
	}

	// Print directory contents as CSV
	if info.Type == "directory" && len(info.Contents) > 0 {
		head := []string{"Name", "Type", "Size", "Extension"}
		w.Write(head)
		for _, c := range info.Contents {
			w.Write([]string{c.Name, string(c.Type), fmt.Sprint(c.Size), c.Extension})
		}
		return
	}

	// Fallback: print all fields as key-value pairs
	w.Write([]string{"Field", "Value"})
	w.Write([]string{"Name", info.Name})
	w.Write([]string{"Type", string(info.Type)})
	w.Write([]string{"Size", fmt.Sprint(info.Size)})
	w.Write([]string{"Extension", info.Extension})
	w.Write([]string{"Content", info.Content})
}
