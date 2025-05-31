package text

import (
	"io/ioutil"
	"lazybox/internal/ir"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
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
		// It's okay if content can't be read (e.g., binary, permissions),
		// still return FileInfo if Lstat succeeded.
		// The error can be stored in FileInfo.Error if desired.
	}
	content := string(contentBytes)

	fileInfo := &ir.FileInfo{
		Name:         filepath.Base(absPath),
		Path:         path, // Use original path for consistency if preferred
		AbsolutePath: absPath,
		Type:         ir.FileTypeFile, // Assume file, can be refined
		Size:         fi.Size(),
		Mode:         fi.Mode().String(),
		ModTime:      fi.ModTime(),
		Extension:    filepath.Ext(absPath),
	}

	if err == nil { // If content was read successfully
		fileInfo.SetContent(content) // Use helper to set *string
		textAnalysis := &ir.TextInfo{
			LineCount: len(strings.Split(content, "\\n")),
			WordCount: len(strings.Fields(content)),
			CharCount: len(content),
			// Other analyses can be added here
		}
		fileInfo.TextAnalysis = textAnalysis
	} else {
		fileInfo.Error = err.Error()
	}

	// Determine if it's a directory
	if fi.IsDir() {
		fileInfo.Type = ir.FileTypeDir
		fileInfo.Content = nil // Directories don't have string content this way
		fileInfo.TextAnalysis = nil
	}

	return fileInfo, nil
}

// Analyze analyzes the given text content and returns detailed TextInfo.
func Analyze(content string, filePath string /* optional, for context */) (*ir.TextInfo, error) {
	textInfo := &ir.TextInfo{
		LineCount:   len(strings.Split(content, "\\n")),
		WordCount:   len(strings.Fields(content)),
		CharCount:   len(content),
		Readability: &ir.ReadabilityScores{},
		Sentiment:   &ir.SentimentAnalysis{}, // Corrected type
		// Keywords will be populated below
		// Encoding and Language detection would require more sophisticated libraries
	}

	// Basic keyword frequency
	words := strings.Fields(strings.ToLower(content))
	wordCounts := make(map[string]int)
	ignoredWords := map[string]bool{
		"the": true, "a": true, "an": true, "is": true, "are": true, "was": true, "were": true,
		"to": true, "of": true, "in": true, "on": true, "at": true, "for": true, "and": true,
		"it": true, "this": true, "that": true, "with": true, "by": true, "as": true,
	}

	nonAlphaNumeric := regexp.MustCompile("[^a-z0-9 ]+")
	for _, word := range words {
		cleanedWord := nonAlphaNumeric.ReplaceAllString(word, "")
		if len(cleanedWord) > 2 && !ignoredWords[cleanedWord] {
			wordCounts[cleanedWord]++
		}
	}

	// Convert map to slice of KeywordFrequency and sort by frequency
	kwList := make([]ir.KeywordFrequency, 0, len(wordCounts))
	for k, v := range wordCounts {
		kwList = append(kwList, ir.KeywordFrequency{Keyword: k, Count: v}) // Corrected field to Count
	}
	sort.Slice(kwList, func(i, j int) bool {
		return kwList[i].Count > kwList[j].Count // Corrected field to Count
	})

	// Store top N keywords (e.g., top 10)
	topN := 10
	if len(kwList) < topN {
		topN = len(kwList)
	}
	textInfo.Keywords = kwList[:topN]

	// TODO: Implement Readability Scores (e.g., Flesch-Kincaid, Gunning Fog)
	// textInfo.Readability.FleschKincaidGradeLevel = ...

	// TODO: Implement Sentiment Analysis
	// textInfo.Sentiment.Polarity = ...
	// textInfo.Sentiment.Subjectivity = ...

	// If filePath is provided, it could be used to enrich TextInfo if needed,
	// but TextInfo primarily focuses on the content itself.
	// If FileInfo is needed alongside TextInfo, they should be associated.
	// For example, a FileInfo struct could have a field for TextAnalysis *TextInfo.

	return textInfo, nil
}
