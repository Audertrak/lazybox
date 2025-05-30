package output

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// PrintFastfetchStyle prints a system info block inspired by fastfetch/neofetch
func PrintFastfetchStyle() {
	osName := runtime.GOOS
	arch := runtime.GOARCH
	user := "?"
	host := "?"
	if u := getEnv("USERNAME"); u != "" {
		user = u
	}
	if h := getEnv("COMPUTERNAME"); h != "" {
		host = h
	}
	goVersion := runtime.Version()

	// Use base16 theme colors
	colors := []lipgloss.Color{
		themeColor(CurrentTheme.Base08),
		themeColor(CurrentTheme.Base09),
		themeColor(CurrentTheme.Base0A),
		themeColor(CurrentTheme.Base0B),
		themeColor(CurrentTheme.Base0C),
		themeColor(CurrentTheme.Base0D),
		themeColor(CurrentTheme.Base0E),
		themeColor(CurrentTheme.Base0F),
	}
	var colorBlocks []string
	for _, c := range colors {
		colorBlocks = append(colorBlocks, lipgloss.NewStyle().Background(c).Foreground(c).Render("   "))
	}
	colorLine := lipgloss.JoinHorizontal(lipgloss.Top, colorBlocks...)

	banner := lipgloss.NewStyle().Bold(true).Foreground(themeColor(CurrentTheme.Base0D)).Render("lazybox")
	infoLines := []string{
		fmt.Sprintf("%s@%s", user, host),
		fmt.Sprintf("OS:   %s", osName),
		fmt.Sprintf("Arch: %s", arch),
		fmt.Sprintf("Go:   %s", goVersion),
	}
	labelStyle := lipgloss.NewStyle().Foreground(themeColor(CurrentTheme.Base0B)).Bold(true)
	valueStyle := lipgloss.NewStyle().Foreground(themeColor(CurrentTheme.Base05))
	var infoStyled []string
	for _, line := range infoLines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			infoStyled = append(infoStyled, labelStyle.Render(parts[0]+":")+" "+valueStyle.Render(strings.TrimSpace(parts[1])))
		} else {
			infoStyled = append(infoStyled, valueStyle.Render(line))
		}
	}
	infoBlock := lipgloss.JoinVertical(lipgloss.Left, infoStyled...)

	block := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(themeColor(CurrentTheme.Base0C)).Padding(1, 4).Background(themeColor(CurrentTheme.Base00)).Render(
		colorLine + "\n" + banner + "\n" + infoBlock + "\n" + colorLine,
	)
	fmt.Println(block)

	// Try to print a logo with figlet if available
	if _, err := exec.LookPath("figlet"); err == nil {
		cmd := exec.Command("figlet", "lazybox")
		out, err := cmd.Output()
		if err == nil {
			fmt.Print(lipgloss.NewStyle().Foreground(themeColor(CurrentTheme.Base0D)).Render(string(out)))
		}
	}
}

func getEnv(key string) string {
	if v, ok := syscallEnv(key); ok {
		return v
	}
	return ""
}

// syscallEnv is a helper for getting env vars cross-platform
func syscallEnv(key string) (string, bool) {
	for _, e := range strings.Split(strings.Join([]string{getEnvWin(), getEnvUnix()}, ";"), ";") {
		if kv := strings.SplitN(e, "=", 2); len(kv) == 2 && kv[0] == key {
			return kv[1], true
		}
	}
	return "", false
}

func getEnvWin() string {
	return strings.Join([]string{"USERNAME=" + os.Getenv("USERNAME"), "COMPUTERNAME=" + os.Getenv("COMPUTERNAME")}, ";")
}

func getEnvUnix() string {
	return strings.Join([]string{"USER=" + os.Getenv("USER"), "HOSTNAME=" + os.Getenv("HOSTNAME")}, ";")
}
