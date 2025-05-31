package file

import (
	"io/ioutil"
	"lazybox/internal/ir"
	"os"
	"path/filepath"
	"strings"
)

// Read reads a file and returns its metadata and contents as an ir.FileInfo
// It will also perform a basic text analysis if the file seems to be text-based.
func Read(path string) (*ir.FileInfo, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	fi, err := os.Lstat(absPath)
	if err != nil {
		return &ir.FileInfo{
			Name:         filepath.Base(absPath),
			Path:         absPath,
			AbsolutePath: absPath,
			Error:        err.Error(),
		}, nil
	}

	fileInfo := &ir.FileInfo{
		Name:         filepath.Base(absPath),
		Path:         absPath,
		AbsolutePath: absPath,
		Type:         ir.FileTypeFile,
		Size:         fi.Size(),
		Mode:         fi.Mode().String(),
		ModTime:      fi.ModTime(),
		Extension:    filepath.Ext(absPath),
	}

	// Attempt to read content
	contentBytes, readErr := ioutil.ReadFile(absPath)
	if readErr == nil {
		contentStr := string(contentBytes) // Create a string variable
		fileInfo.Content = &contentStr     // Assign its address

		// Initialize TextAnalysis if it's nil
		if fileInfo.TextAnalysis == nil {
			fileInfo.TextAnalysis = &ir.TextInfo{}
		}
		// Populate TextAnalysis fields
		if fileInfo.Content != nil { // Check if content is not nil before dereferencing
			fileInfo.TextAnalysis.LineCount = len(strings.Split(*fileInfo.Content, "\n"))
			fileInfo.TextAnalysis.WordCount = len(strings.Fields(*fileInfo.Content))
		} else {
			fileInfo.TextAnalysis.LineCount = 0
			fileInfo.TextAnalysis.WordCount = 0
		}
		// TODO: Add a heuristic to decide if it's a text file suitable for deeper analysis
		// For now, we assume all readable files might be text.
		// Potentially, we could create a *ir.TextInfo here and embed it or link it.
		// However, the current `file` command expects `*ir.FileInfo`.
		// We can enhance `handleOutput` or the command itself to expect `*ir.TextInfo`
		// if more detailed text analysis is always desired for the `file` target.
	} else {
		fileInfo.Error = readErr.Error() // Store read error if any
	}

	// Populate OS-specific fields (placeholder)
	// fileInfo.Owner = getOwner(absPath) // To be implemented
	// fileInfo.Group = getGroup(absPath) // To be implemented
	// fileInfo.CreateTime = getCreateTime(absPath) // To be implemented

	return fileInfo, nil
}

// Helper functions for OS-specific info (to be implemented)
// func getOwner(path string) string { /* ... */ return "" }
// func getGroup(path string) string { /* ... */ return "" }
// func getCreateTime(path string) *time.Time { /* ... */ return nil }
