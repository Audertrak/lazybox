package pkg

import (
	"lazybox/internal/fs"
	"lazybox/internal/ir"
)

// Crawl is a placeholder for package crawling. For now, just uses fs.Scan.
func Crawl(path string) (*ir.FileInfo, error) {
	return fs.Scan(path)
}
