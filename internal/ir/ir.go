package ir

import (
	"os"
	"path/filepath"
	"time"
)

// FileType is an enumeration for different types of files or data.
type FileType string

const (
	FileTypeFile         FileType = "file"
	FileTypeDir          FileType = "directory"
	FileTypeSymlink      FileType = "symlink"
	FileTypeCode         FileType = "code"
	FileTypeAPI          FileType = "api"
	FileTypeEnum         FileType = "enum"
	FileTypeStruct       FileType = "struct"
	FileTypePackage      FileType = "package"
	FileTypeList         FileType = "list"
	FileTypeText         FileType = "text"
	FileTypeUnknown      FileType = "unknown"
	FileTypeGitRepo      FileType = "git_repository"
	FileTypeGitSubmodule FileType = "git_submodule"
)

// FileInfo represents the intermediate representation for a file or directory.
type FileInfo struct {
	Name             string                 `json:"name"`
	Path             string                 `json:"path"`
	AbsolutePath     string                 `json:"absolute_path"`
	Type             FileType               `json:"type"`
	IsDir            bool                   `json:"is_dir"`
	IsSymlink        bool                   `json:"is_symlink,omitempty"`
	SymlinkTarget    string                 `json:"symlink_target,omitempty"`
	Size             int64                  `json:"size"`
	Mode             string                 `json:"mode,omitempty"` // Store as string from os.FileMode.String()
	ModTime          time.Time              `json:"mod_time,omitempty"`
	CreateTime       time.Time              `json:"create_time,omitempty"` // OS-dependent
	Owner            string                 `json:"owner,omitempty"`       // OS-dependent
	Group            string                 `json:"group,omitempty"`       // OS-dependent
	Extension        string                 `json:"extension,omitempty"`
	Content          *string                `json:"content,omitempty"` // Pointer to allow nil for non-text or large files not read
	Error            string                 `json:"error,omitempty"`
	Children         []*FileInfo            `json:"children,omitempty"`           // For directories
	TextAnalysis     *TextInfo              `json:"text_analysis,omitempty"`      // For text files
	GitRemoteURL     string                 `json:"git_remote_url,omitempty"`     // For git repositories
	GitCurrentBranch string                 `json:"git_current_branch,omitempty"` // For git repositories
	Metadata         map[string]interface{} `json:"metadata,omitempty"`           // For any other specific metadata
}

// KeywordFrequency stores a keyword and its count.
type KeywordFrequency struct {
	Keyword string `json:"keyword"`
	Count   int    `json:"count"`
}

// ReadabilityScores stores various readability metrics.
type ReadabilityScores struct {
	FleschKincaidGradeLevel float64 `json:"flesch_kincaid_grade_level,omitempty"`
	GunningFogIndex         float64 `json:"gunning_fog_index,omitempty"`
	// Add other scores as needed
}

// SentimentAnalysis stores sentiment scores.
type SentimentAnalysis struct {
	Polarity     float64 `json:"polarity,omitempty"`     // e.g., -1 (negative) to 1 (positive)
	Subjectivity float64 `json:"subjectivity,omitempty"` // e.g., 0 (objective) to 1 (subjective)
}

// TextInfo holds detailed analysis of text content.
type TextInfo struct {
	LineCount             int                `json:"line_count"`
	WordCount             int                `json:"word_count"`
	CharCount             int                `json:"char_count"`
	Keywords              []KeywordFrequency `json:"keywords,omitempty"`
	DetectedLanguage      string             `json:"detected_language,omitempty"`
	Readability           *ReadabilityScores `json:"readability,omitempty"`
	Sentiment             *SentimentAnalysis `json:"sentiment,omitempty"`
	IsBinary              bool               `json:"is_binary,omitempty"`
	MimeType              string             `json:"mime_type,omitempty"`
	Encoding              string             `json:"encoding,omitempty"`
	AverageWordLength     float64            `json:"average_word_length,omitempty"`
	AverageSentenceLength float64            `json:"average_sentence_length,omitempty"`
}

// NewFileInfo creates a basic FileInfo struct.
func NewFileInfo(name, path, absPath string, fileType FileType, isDir bool, size int64, mode os.FileMode, modTime time.Time) *FileInfo {
	fi := &FileInfo{
		Name:         name,
		Path:         path,
		AbsolutePath: absPath,
		Type:         fileType,
		IsDir:        isDir,
		Size:         size,
		ModTime:      modTime,
		Metadata:     make(map[string]interface{}),
	}
	if mode != 0 {
		fi.Mode = mode.String()
	}
	if !isDir {
		// Safely get extension, ensuring name is not empty and contains a dot.
		if ext := filepath.Ext(name); ext != "" && len(name) > len(ext) {
			fi.Extension = ext
		}
	}
	return fi
}

// Helper function to safely set content
func (fi *FileInfo) SetContent(content string) {
	fi.Content = &content
}

// AddChild adds a child FileInfo to the current FileInfo.
func (fi *FileInfo) AddChild(child *FileInfo) {
	fi.Children = append(fi.Children, child)
}
