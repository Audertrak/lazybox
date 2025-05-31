package fs

import (
	"fmt"
	"io/fs"
	"lazybox/internal/ir"
	"os"
	"path/filepath"
	"strings"
	"time"
	// "syscall" // For owner/group/createtime - OS specific, handle later
)

// Scan recursively scans a directory or gets info for a file.
func Scan(path string) (*ir.FileInfo, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for %s: %w", path, err)
	}

	info, err := os.Lstat(absPath) // Use Lstat to get info about symlink itself
	if err != nil {
		return nil, fmt.Errorf("failed to stat %s: %w", absPath, err)
	}

	fileIR := &ir.FileInfo{
		Name:         info.Name(),
		Path:         path, // Original path provided
		AbsolutePath: absPath,
		Size:         info.Size(),
		Mode:         info.Mode().String(),
		ModTime:      info.ModTime(),
		Extension:    strings.ToLower(filepath.Ext(info.Name())),
	}

	// Type determination
	if info.IsDir() {
		fileIR.Type = ir.FileTypeDir // Corrected: FileTypeDir
	} else if info.Mode()&os.ModeSymlink != 0 {
		fileIR.Type = ir.FileTypeSymlink
		symlinkTarget, err := os.Readlink(absPath)
		if err != nil {
			fileIR.Error = fmt.Sprintf("failed to read symlink target: %v", err)
		} else {
			fileIR.SymlinkTarget = symlinkTarget
		}
	} else {
		fileIR.Type = ir.FileTypeFile
	}

	// Owner, Group, CreateTime - OS specific, placeholder for now
	// ct := getCreationTime(info) // Placeholder for platform-specific logic
	// if !ct.IsZero() {
	// 	fileIR.CreateTime = &ct
	// }

	// Git repository detection
	isGit, gitDir := isGitRepo(absPath)
	if fileIR.Metadata == nil {
		fileIR.Metadata = make(map[string]interface{})
	}
	fileIR.Metadata["is_git_repo"] = isGit // Corrected: Use Metadata
	if isGit {
		remotes, err := getGitRemotes(gitDir)
		if err == nil {
			fileIR.Metadata["git_remotes"] = remotes // Corrected: Use Metadata
		} else {
			// Optionally log this error or store it in fileIR.Error
			// fileIR.Error += fmt.Sprintf(\"; git remote error: %v\", err)
		}
	}

	if fileIR.Type == ir.FileTypeDir { // Corrected: FileTypeDir
		entries, err := os.ReadDir(absPath)
		if err != nil {
			fileIR.Error += fmt.Sprintf("; failed to read directory %s: %v", absPath, err)
			return fileIR, nil
		}

		fileIR.Children = make([]*ir.FileInfo, 0, len(entries)) // Corrected: Children
		for _, entry := range entries {
			entryPath := filepath.Join(absPath, entry.Name())

			entryIR, err := Scan(entryPath) // Recursive call
			if err != nil {
				errorEntryIR := &ir.FileInfo{
					Name:         entry.Name(),
					Path:         filepath.Join(path, entry.Name()),
					AbsolutePath: entryPath,
					Error:        err.Error(),
				}
				entryInfo, statErr := os.Lstat(entryPath)
				if statErr == nil {
					if entryInfo.IsDir() {
						errorEntryIR.Type = ir.FileTypeDir // Corrected: FileTypeDir
					} else if entryInfo.Mode()&os.ModeSymlink != 0 {
						errorEntryIR.Type = ir.FileTypeSymlink
					} else {
						errorEntryIR.Type = ir.FileTypeFile
					}
				}
				fileIR.Children = append(fileIR.Children, errorEntryIR) // Corrected: Children
				continue
			}
			fileIR.Children = append(fileIR.Children, entryIR) // Corrected: Children
		}
	}

	return fileIR, nil
}

// isGitRepo checks if the given path or any of its parents is a Git repository.
// It returns true and the path to the .git directory if found.
func isGitRepo(path string) (bool, string) {
	currentPath := path
	if fi, err := os.Stat(path); err == nil && !fi.IsDir() { // If path is a file, start check from its directory
		currentPath = filepath.Dir(path)
	}

	for {
		gitPath := filepath.Join(currentPath, ".git")
		stat, err := os.Stat(gitPath)
		if err == nil && stat.IsDir() {
			return true, gitPath
		}

		parent := filepath.Dir(currentPath)
		if parent == currentPath { // Reached root or invalid path
			break
		}
		currentPath = parent
	}
	return false, ""
}

// getGitRemotes parses the .git/config file to find remote origins.
func getGitRemotes(gitDir string) (map[string]string, error) {
	configPath := filepath.Join(filepath.Dir(gitDir), ".git", "config") // Ensure it's .git/config
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not read git config %s: %w", configPath, err)
	}

	remotes := make(map[string]string)
	lines := strings.Split(string(content), "\n")
	var currentRemoteName string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "[remote \"") && strings.HasSuffix(trimmedLine, "\"]") {
			name := strings.TrimPrefix(trimmedLine, "[remote \"")
			name = strings.TrimSuffix(name, "\"]")
			currentRemoteName = name
		} else if currentRemoteName != "" && strings.HasPrefix(trimmedLine, "url = ") {
			url := strings.TrimPrefix(trimmedLine, "url = ")
			remotes[currentRemoteName] = url
			currentRemoteName = ""
		} else if strings.HasPrefix(trimmedLine, "[") && currentRemoteName != "" {
			// New section started before URL for current remote was found
			currentRemoteName = ""
		}
	}
	// No error if remotes is empty, it's a valid state.
	return remotes, nil
}

// Placeholder for getCreationTime - requires OS-specific implementation
func getCreationTime(info fs.FileInfo) time.Time {
	// This is highly OS-dependent.
	// On Linux, info.Sys().(*syscall.Stat_t).Ctim is change time.
	// On macOS/BSDs, info.Sys().(*syscall.Stat_t).Birthtimespec.
	// On Windows, info.Sys().(*syscall.Win32FileAttributeData).CreationTime.
	// For now, returning zero time.
	return time.Time{}
}

// Platform-specific implementations are in fs_unix.go and fs_windows.go
