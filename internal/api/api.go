package api

import (
	"io/ioutil"
	"lazybox/internal/ir"
	"os"
	"path/filepath"
)

// Extract is a placeholder for API extraction. For now, just reads file content and metadata.
func Extract(path string) (*ir.FileInfo, error) {
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
	contentBytes, err := ioutil.ReadFile(absPath)
	content := ""
	if err == nil {
		content = string(contentBytes)
	}
	return &ir.FileInfo{
		Name:         filepath.Base(absPath),
		Path:         absPath,
		AbsolutePath: absPath,
		Type:         ir.FileTypeFile,
		Size:         fi.Size(),
		Mode:         fi.Mode().String(),
		ModTime:      fi.ModTime(),
		Extension:    filepath.Ext(absPath),
		Content:      content,
	}, nil
}
