package output

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
