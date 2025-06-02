package output

import (
	"encoding/json"
	"fmt"
	"lazybox/internal/glpg"
	"lazybox/internal/ir"
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

// RustStyleEntry matches the Rust JSON output for fs target
// Only fields present in the Rust output, camelCase, omitempty
// Recursively used for contents
type RustStyleEntry struct {
	Name             string            `json:"name"`
	Path             *string           `json:"path,omitempty"`
	AbsolutePathFull *string           `json:"absolutePathFull,omitempty"`
	EntryTypeFull    *string           `json:"entryTypeFull,omitempty"`
	IsSymlink        *bool             `json:"isSymlink,omitempty"`
	SymlinkTarget    *string           `json:"symlinkTarget,omitempty"`
	IsGitRepo        *bool             `json:"isGitRepo,omitempty"`
	GitRemotes       map[string]string `json:"gitRemotes,omitempty"`
	Contents         []RustStyleEntry  `json:"contents,omitempty"`
	Error            *string           `json:"error,omitempty"`
}

// Convert from *ir.FileInfo to RustStyleEntry recursively
func FileInfoToRustStyleEntry(fi *ir.FileInfo, compact bool, isRoot bool, rootPath string) RustStyleEntry {
	var path *string
	if compact {
		if isRoot {
			p := rootPath
			path = &p
		} else {
			rel := fi.Path
			if rel == rootPath {
				rel = fi.Name
			} else if rel != "" && len(rootPath) > 0 && len(fi.Path) > len(rootPath) && fi.Path[:len(rootPath)] == rootPath {
				rel = fi.Path[len(rootPath):]
				if len(rel) > 0 && (rel[0] == '/' || rel[0] == '\\') {
					rel = rel[1:]
				}
			}
			path = &rel
		}
	} else {
		if isRoot {
			ap := fi.AbsolutePath
			path = &ap
		}
	}
	var contents []RustStyleEntry
	for _, child := range fi.Children {
		contents = append(contents, FileInfoToRustStyleEntry(child, compact, false, rootPath))
	}
	if compact {
		return RustStyleEntry{
			Name:     fi.Name,
			Path:     path,
			Contents: contents,
		}
	}
	// Normal output: all fields
	var absolutePathFull *string
	var entryTypeFull *string
	var isSymlink *bool
	var symlinkTarget *string
	var isGitRepo *bool
	var gitRemotes map[string]string
	var errorStr *string
	if isRoot {
		ap := fi.AbsolutePath
		absolutePathFull = &ap
		et := string(fi.Type)
		entryTypeFull = &et
	}
	if fi.IsSymlink {
		b := true
		isSymlink = &b
	} else if fi.Type == ir.FileTypeSymlink {
		b := true
		isSymlink = &b
	}
	if fi.SymlinkTarget != "" {
		st := fi.SymlinkTarget
		symlinkTarget = &st
	}
	if v, ok := fi.Metadata["is_git_repo"]; ok {
		if b, ok2 := v.(bool); ok2 {
			isGitRepo = &b
		}
	}
	if v, ok := fi.Metadata["git_remotes"]; ok {
		if m, ok2 := v.(map[string]string); ok2 {
			gitRemotes = m
		}
	}
	if fi.Error != "" {
		e := fi.Error
		errorStr = &e
	}
	return RustStyleEntry{
		Name:             fi.Name,
		Path:             path,
		AbsolutePathFull: absolutePathFull,
		EntryTypeFull:    entryTypeFull,
		IsSymlink:        isSymlink,
		SymlinkTarget:    symlinkTarget,
		IsGitRepo:        isGitRepo,
		GitRemotes:       gitRemotes,
		Contents:         contents,
		Error:            errorStr,
	}
}

// PrintJSONWithHighlight prints JSON with syntax highlighting using lipgloss
func PrintJSONWithHighlight(data []byte, minified bool) {
	var out any
	if err := json.Unmarshal(data, &out); err != nil {
		fmt.Println(string(data))
		return
	}
	printJSONValueWithHighlight(out, 0, minified)
}

func printJSONValueWithHighlight(v any, indent int, minified bool) {
	pad := func(n int) string {
		if minified {
			return ""
		}
		return string(make([]byte, n*2))
	}
	switch val := v.(type) {
	case map[string]any:
		if minified {
			fmt.Print(jqKeyStyle.Render("{"))
		} else {
			fmt.Print(jqKeyStyle.Render("{\n"))
		}
		first := true
		for k, v2 := range val {
			if !first {
				if minified {
					fmt.Print(jqKeyStyle.Render(","))
				} else {
					fmt.Print(jqKeyStyle.Render(",\n"))
				}
			}
			first = false
			fmt.Print(pad(indent + 1))
			fmt.Print(jqKeyStyle.Render("\"" + k + "\""))
			fmt.Print(jqKeyStyle.Render(": "))
			printJSONValueWithHighlight(v2, indent+1, minified)
		}
		if minified {
			fmt.Print(jqKeyStyle.Render("}"))
		} else {
			fmt.Print("\n" + pad(indent) + jqKeyStyle.Render("}"))
		}
	case []any:
		fmt.Print(jqKeyStyle.Render("["))
		for i, v2 := range val {
			if i > 0 {
				if minified {
					fmt.Print(jqKeyStyle.Render(","))
				} else {
					fmt.Print(jqKeyStyle.Render(", "))
				}
			}
			printJSONValueWithHighlight(v2, indent+1, minified)
		}
		fmt.Print(jqKeyStyle.Render("]"))
	case string:
		fmt.Print(jqStringStyle.Render("\"" + val + "\""))
	case float64:
		fmt.Print(jqNumStyle.Render(fmt.Sprintf("%v", val)))
	case bool:
		fmt.Print(jqBoolStyle.Render(fmt.Sprintf("%v", val)))
	case nil:
		fmt.Print(jqNullStyle.Render("null"))
	default:
		fmt.Print(jqStringStyle.Render(fmt.Sprintf("%v", val)))
	}
}

// PrintGLPGAsJSON serializes the GLPG to a Rust-style recursive tree for fs target, with pretty/minified and normal/compact modes
func PrintGLPGAsJSON(graph *glpg.GLPG, flags map[string]bool) error {
	compact := flags["less"] || flags["compact"]
	minified := flags["min"]

	if graph.OriginalFileInfo != nil {
		rootPath := graph.OriginalFileInfo.AbsolutePath
		entry := FileInfoToRustStyleEntry(graph.OriginalFileInfo, compact, true, rootPath)
		var data []byte
		var err error
		if minified {
			data, err = json.Marshal(entry)
		} else {
			data, err = json.MarshalIndent(entry, "", "  ")
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling RustStyleEntry: %v\n", err)
			return err
		}
		if minified {
			PrintJSONWithHighlight(data, true)
		} else {
			PrintJSONWithHighlight(data, false)
			fmt.Println()
		}
		return nil
	}

	// fallback: original GLPG graph output for non-fs targets
	var data []byte
	var err error
	if minified {
		data, err = json.Marshal(graph)
	} else {
		data, err = json.MarshalIndent(graph, "", "  ")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling GLPG to JSON: %v\n", err)
		return err
	}
	PrintJSONWithHighlight(data, minified)
	if !minified {
		fmt.Println()
	}
	return nil
}

// PrintGLPGAsMinJSON is now an alias for PrintGLPGAsJSON with minified flag
func PrintGLPGAsMinJSON(graph *glpg.GLPG) error {
	return PrintGLPGAsJSON(graph, map[string]bool{"min": true})
}

// PrintGLPGAsLessJSON is now an alias for PrintGLPGAsJSON with compact flag
func PrintGLPGAsLessJSON(graph *glpg.GLPG) error {
	return PrintGLPGAsJSON(graph, map[string]bool{"less": true})
}
