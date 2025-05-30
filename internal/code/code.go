package code

import (
	"io/ioutil"
	"lazybox/internal/ir"
	"os"
	"path/filepath"
)

func ExtractPlaceholder(path string) (*ir.FileInfo, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	fi, err := os.Lstat(absPath)
	if err != nil {
		return nil, err
	}
	contentBytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	content := string(contentBytes)
	return &ir.FileInfo{
		Name:      filepath.Base(absPath),
		Path:      absPath,
		Type:      "code",
		Size:      fi.Size(),
		Content:   content,
		Extension: filepath.Ext(absPath),
		Mode:      fi.Mode().String(),
	}, nil
}
