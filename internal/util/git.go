package util

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func IsGitRepo(path string) bool {
	gitDir := filepath.Join(path, ".git")
	fi, err := os.Stat(gitDir)
	return err == nil && fi.IsDir()
}

func GetGitRemotes(path string) map[string]string {
	remotes := make(map[string]string)
	cmd := exec.Command("git", "remote", "-v")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return remotes
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			name := fields[0]
			url := fields[1]
			if _, exists := remotes[name]; !exists {
				remotes[name] = url
			}
		}
	}
	return remotes
}
