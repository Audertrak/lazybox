package output

import (
	"fmt"
	"lazybox/internal/theme" // Import theme package

	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure" // Import go-figure
)

// PrintBannerText prints the application title as a FIGlet-style banner
// using the go-figure library and styled with lipgloss according to the current theme.
func PrintBannerText(text string) {
	ct := theme.GetDefaultTheme() // Get the currently loaded theme

	// Create a FIGlet figure.
	// "standard" is a common FIGlet font. Explore others if desired.
	myFigure := figure.NewFigure(text, "standard", true)

	// Style the FIGlet text using lipgloss.
	// Using a prominent color from the theme for the banner text.
	bannerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ct.Base0E)). // Example: A bright, attention-grabbing color from the theme
		Bold(true).
		MarginBottom(1)

	// Print each line of the FIGlet text with the style applied.
	// myFigure.Slicify() returns a slice of strings, each being a line of the FIGlet text.
	for _, line := range myFigure.Slicify() {
		fmt.Println(bannerStyle.Render(line))
	}

	// Optional: Add a tagline or version number below the banner, also styled.
	taglineStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ct.Base0B)). // A complementary color from the theme
		Italic(true).
		MarginBottom(1) // Add some space below the tagline

	// You can make the tagline more dynamic, e.g., include version info if available.
	fmt.Println(taglineStyle.Render("  Your polymorphic structured data swiss army knife..."))
	fmt.Println() // Add an extra newline for better spacing before subsequent output
}

/*
// Old implementation using external figlet/boxes commands (commented out)
import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Try to use figlet or boxes if available, else fallback to lipgloss styled text
func PrintBannerText(text string) {
	// Try figlet first
	if _, err := exec.LookPath("figlet"); err == nil {
		cmd := exec.Command("figlet", text)
		out, err := cmd.Output()
		if err == nil {
			fmt.Print(lipgloss.NewStyle().Foreground(themeColor(CurrentTheme.Base0D)).Render(string(out)))
			return
		}
	}
	// Try boxes if available
	if _, err := exec.LookPath("boxes"); err == nil {
		cmd := exec.Command("boxes")
		cmd.Stdin = strings.NewReader(text)
		out, err := cmd.Output()
		if err == nil {
			fmt.Print(lipgloss.NewStyle().Foreground(themeColor(CurrentTheme.Base0C)).Render(string(out)))
			return
		}
	}
	// Fallback: lipgloss styled text
	banner := lipgloss.NewStyle().Bold(true).Foreground(themeColor(CurrentTheme.Base0D)).Background(themeColor(CurrentTheme.Base01)).Padding(1, 4).Render(text)
	fmt.Println(banner)
}
*/
