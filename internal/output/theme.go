package output

import "github.com/charmbracelet/lipgloss"

// Base16 color scheme struct
// See: https://github.com/chriskempson/base16/blob/main/styling.md
// Rose Pine default values from https://rosepinetheme.com/palette/
type Base16Theme struct {
	Base00 string // background
	Base01 string // lighter background
	Base02 string // selection background
	Base03 string // comments, invisibles, line highlight
	Base04 string // dark foreground
	Base05 string // default foreground
	Base06 string // light foreground
	Base07 string // light background
	Base08 string // variables, tags, deleted
	Base09 string // integers, booleans, constants
	Base0A string // classes, bold, search
	Base0B string // strings, code, inserted
	Base0C string // support, regex, quotes
	Base0D string // functions, methods, headings
	Base0E string // keywords, storage, italic
	Base0F string // deprecated, embedded
	Accent string // accent color
}

// Rose Pine Bloom (default)
var RosePineTheme = Base16Theme{
	Base00: "191724",
	Base01: "1f1d2e",
	Base02: "26233a",
	Base03: "6e6a86",
	Base04: "908caa",
	Base05: "e0def4",
	Base06: "e0def4",
	Base07: "524f67",
	Base08: "eb6f92",
	Base09: "f6c177",
	Base0A: "ebbcba",
	Base0B: "31748f",
	Base0C: "9ccfd8",
	Base0D: "c4a7e7",
	Base0E: "f6c177",
	Base0F: "524f67",
	Accent: "9ccfd8",
}

// Global theme (can be swapped)
var CurrentTheme = RosePineTheme

// Helper to get a lipgloss.Color from a theme hex
func themeColor(hex string) lipgloss.Color {
	if len(hex) == 6 {
		return lipgloss.Color("#" + hex)
	}
	return lipgloss.Color(hex)
}
