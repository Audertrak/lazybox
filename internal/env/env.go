package env

import "lazybox/internal/ir"

func ExtractPlaceholder() (*ir.FileInfo, error) {
	return &ir.FileInfo{
		Name:    "env",
		Type:    "env",
		Content: "[env extraction not yet implemented]",
	}, nil
}
