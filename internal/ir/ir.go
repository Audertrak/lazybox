package ir

import (
	"time"
)

type FileType string

const (
	FileTypeFile      FileType = "file"
	FileTypeDirectory FileType = "directory"
	FileTypeSymlink   FileType = "symlink"
	FileTypeOther     FileType = "other"
)

type FileInfo struct {
	Name          string            `json:"name"`
	Path          string            `json:"path"`
	AbsolutePath  string            `json:"absolute_path"`
	Type          FileType          `json:"type"`
	Size          int64             `json:"size_bytes"`
	Mode          string            `json:"mode"`
	Owner         string            `json:"owner,omitempty"`
	Group         string            `json:"group,omitempty"`
	ModTime       time.Time         `json:"last_modified"`
	CreateTime    *time.Time        `json:"created,omitempty"`
	SymlinkTarget string            `json:"symlink_target,omitempty"`
	IsGitRepo     bool              `json:"is_git_repo,omitempty"`
	GitRemotes    map[string]string `json:"git_remotes,omitempty"`
	Extension     string            `json:"extension,omitempty"`
	Contents      []*FileInfo       `json:"contents,omitempty"`
	Error         string            `json:"error,omitempty"`
}
