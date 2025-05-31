package theme

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Base16Color represents a hex color string.
type Base16Color string

// Base16Theme defines the structure for a Base16 theme.
// It includes common names and the 16 base colors.
type Base16Theme struct {
	Scheme string      `yaml:"scheme"`
	Author string      `yaml:"author"`
	Base00 Base16Color `yaml:"base00"` // Default Background
	Base01 Base16Color `yaml:"base01"` // Lighter Background (Status Bar, Line Numbers)
	Base02 Base16Color `yaml:"base02"` // Selection Background
	Base03 Base16Color `yaml:"base03"` // Comments, Invisibles, Line Highlighting
	Base04 Base16Color `yaml:"base04"` // Dark Foreground (Status Bar)
	Base05 Base16Color `yaml:"base05"` // Default Foreground
	Base06 Base16Color `yaml:"base06"` // Light Foreground
	Base07 Base16Color `yaml:"base07"` // Light Background
	Base08 Base16Color `yaml:"base08"` // Variables, XML Tags, Markup Link Text, Markup Lists, Diff Deleted
	Base09 Base16Color `yaml:"base09"` // Integers, Boolean, Constants, XML Attributes, Markup Link Url
	Base0A Base16Color `yaml:"base0A"` // Classes, Markup Bold, Search Text Background
	Base0B Base16Color `yaml:"base0B"` // Strings, Inherited Class, Markup Code, Diff Inserted
	Base0C Base16Color `yaml:"base0C"` // Support, Regular Expressions, Escape Characters, Markup Quotes
	Base0D Base16Color `yaml:"base0D"` // Functions, Methods, Attribute IDs, Headings
	Base0E Base16Color `yaml:"base0E"` // Keywords, Storage, Selector, Markup Italic, Diff Changed
	Base0F Base16Color `yaml:"base0F"` // Deprecated, Opening/Closing Embedded Language Tags
}

// DefaultTheme holds the currently active theme.
var DefaultTheme *Base16Theme

// LoadTheme loads a Base16 theme from a YAML file.
func LoadTheme(filePath string) (*Base16Theme, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme file %s: %w", filePath, err)
	}

	var theme Base16Theme
	err = yaml.Unmarshal(data, &theme)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal theme YAML from %s: %w", filePath, err)
	}
	return &theme, nil
}

// SetDefaultTheme sets the global default theme.
func SetDefaultTheme(theme *Base16Theme) {
	DefaultTheme = theme
}

// GetDefaultTheme retrieves the global default theme.
// If no theme is set, it attempts to load Rose Pine as a fallback.
func GetDefaultTheme() *Base16Theme {
	if DefaultTheme == nil {
		// Attempt to load from env or a default path first (optional)
		// For now, directly fallback to Rose Pine
		DefaultTheme = GetRosePineTheme()
	}
	return DefaultTheme
}

// GetRosePineTheme returns the Rose Pine Base16 theme.
func GetRosePineTheme() *Base16Theme {
	return &Base16Theme{
		Scheme: "Rose Pine",
		Author: "Rose Pine",
		Base00: "#191724", // Background
		Base01: "#1f1d2e", // Lighter Background (Used for status bars, line number background)
		Base02: "#26233a", // Selection Background
		Base03: "#555169", // Comments, Invisibles, Line Highlighting
		Base04: "#6e6a86", // Dark Foreground (Used for status bars)
		Base05: "#e0def4", // Default Foreground
		Base06: "#f0f0f0", // Light Foreground (Not often used)
		Base07: "#ffffff", // Light Background (Not often used)
		Base08: "#eb6f92", // Red (Variables, XML Tags, Markup Link Text, Markup Lists, Diff Deleted)
		Base09: "#f6c177", // Orange (Integers, Boolean, Constants, XML Attributes, Markup Link Url)
		Base0A: "#ebbcba", // Yellow (Classes, Markup Bold, Search Text Background)
		Base0B: "#31748f", // Green (Strings, Inherited Class, Markup Code, Diff Inserted)
		Base0C: "#9ccfd8", // Cyan (Support, Regular Expressions, Escape Characters, Markup Quotes)
		Base0D: "#c4a7e7", // Magenta (Functions, Methods, Attribute IDs, Headings)
		Base0E: "#ea9a97", // Violet (Keywords, Storage, Selector, Markup Italic, Diff Changed)
		Base0F: "#524f67", // Brown (Deprecated, Opening/Closing Embedded Language Tags e.g. `<?php ?>`)
	}
}

// Initialize loads the default theme (Rose Pine for now) or a theme specified by an environment variable.
func Initialize() {
	// For now, always initialize with Rose Pine.
	// Later, this could check an environment variable (e.g., LAZYBOX_THEME_PATH)
	// or a configuration file.
	themePath := os.Getenv("LAZYBOX_THEME_PATH")
	if themePath != "" {
		loadedTheme, err := LoadTheme(themePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to load theme from %s: %v. Falling back to default.\n", themePath, err)
			SetDefaultTheme(GetRosePineTheme())
		} else {
			SetDefaultTheme(loadedTheme)
			return
		}
	} else {
		SetDefaultTheme(GetRosePineTheme())
	}
}

// Helper to get a color, falling back to foreground/background if a specific theme color is not set
// This is a basic example; more sophisticated fallback logic might be needed.
// This function might not be strictly necessary if lipgloss.Color can directly take hex strings.
// However, it can be useful for centralizing color access if more complex logic is needed.
func GetColor(baseColor Base16Color) string {
	return string(baseColor)
}
