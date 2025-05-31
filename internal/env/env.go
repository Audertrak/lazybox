package env

import "lazybox/internal/ir"

func ExtractPlaceholder() (*ir.FileInfo, error) {
	content := "[env extraction not yet implemented]" // Store literal in a variable
	return &ir.FileInfo{
		Name:    "env",
		Type:    "env",
		Content: &content, // Assign address of content
	}, nil
}
