package output

import (
	"encoding/json"
	"fmt"
	"lazybox/internal/ir"
	"os"
	"regexp"
	"strings"

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
	jqKeyStyle    = lipgloss.NewStyle().Foreground(themeColor(CurrentTheme.Base0B)).Bold(true)
	jqStringStyle = lipgloss.NewStyle().Foreground(themeColor(CurrentTheme.Base0A))
	jqNumStyle    = lipgloss.NewStyle().Foreground(themeColor(CurrentTheme.Base09))
	jqBoolStyle   = lipgloss.NewStyle().Foreground(themeColor(CurrentTheme.Base0E)).Bold(true)
	jqNullStyle   = lipgloss.NewStyle().Foreground(themeColor(CurrentTheme.Base03)).Italic(true)
)

func PrintJSON(info *ir.FileInfo) {
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(data))
}

func PrintMinJSON(info *ir.FileInfo) {
	data, err := json.Marshal(info)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(data))
}

func PrintLessJSON(info *ir.FileInfo) {
	if info == nil {
		fmt.Println("{}")
		return
	}
	less := map[string]interface{}{
		"name": info.Name,
		"type": info.Type,
		"path": info.Path,
	}
	if info.Type == ir.FileTypeFile {
		less["extension"] = info.Extension
		less["size_bytes"] = info.Size
	}
	if info.Type == ir.FileTypeDirectory {
		less["contents_count"] = len(info.Contents)
	}
	data, err := json.MarshalIndent(less, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(data))
}

func TokenizeText(s string) string {
	s = strings.ToLower(s)
	for k, v := range phraseSubs {
		s = strings.ReplaceAll(s, k, v)
	}
	s = fillerWords.ReplaceAllString(s, "")
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	return s
}

func PrintTokenized(info *ir.FileInfo) {
	if info == nil {
		fmt.Println("")
		return
	}
	fmt.Printf("name:%s type:%s path:%s", info.Name, info.Type, info.Path)
	if info.Type == ir.FileTypeFile {
		fmt.Printf(" extension:%s size:%d", info.Extension, info.Size)
		if info.Content != "" {
			fmt.Print("\ncontent:")
			// Only print first 500 chars for brevity
			short := info.Content
			if len(short) > 500 {
				short = short[:500]
			}
			fmt.Print(TokenizeText(short))
		}
	}
	if info.Type == ir.FileTypeDirectory {
		fmt.Printf(" contents:%d", len(info.Contents))
	}
	fmt.Println()
}

func PrintJQ(info *ir.FileInfo) {
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		styled := line
		styled = regexp.MustCompile(`"([^"]+)":`).ReplaceAllStringFunc(styled, func(m string) string {
			key := m[1 : len(m)-2]
			return jqKeyStyle.Render("\""+key+"\"") + ":"
		})
		styled = regexp.MustCompile(`: "([^"]*)"`).ReplaceAllStringFunc(styled, func(m string) string {
			val := m[3:]
			return ": " + jqStringStyle.Render(val)
		})
		styled = regexp.MustCompile(`: (\d+)`).ReplaceAllStringFunc(styled, func(m string) string {
			val := m[2:]
			return ": " + jqNumStyle.Render(val)
		})
		styled = regexp.MustCompile(`: (true|false)`).ReplaceAllStringFunc(styled, func(m string) string {
			val := m[2:]
			return ": " + jqBoolStyle.Render(val)
		})
		styled = regexp.MustCompile(`: null`).ReplaceAllStringFunc(styled, func(m string) string {
			return ": " + jqNullStyle.Render("null")
		})
		fmt.Println(styled)
	}
}

func Print(info *ir.FileInfo, mode string) {
	switch mode {
	case "jsonify", "json", "":
		PrintJSON(info)
	case "mdify", "markdown":
		PrintMarkdown(info)
	case "prettify", "pretty":
		PrintPretty(info)
	case "tabelify", "table":
		PrintTable(info)
	case "commafy", "csv":
		PrintCSV(info)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown or unsupported mode: %s\n", mode)
		os.Exit(2)
	}
}
