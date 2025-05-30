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
		return nil, err
	}
	content := string(contentBytes)
	lines := strings.Count(content, "\n")
	words := len(strings.Fields(content))
	return &ir.FileInfo{
		Name:      filepath.Base(absPath),
		Path:      absPath,
		Type:      "text",
		Size:      fi.Size(),
		Content:   content,
		Extension: filepath.Ext(absPath),
		Mode:      fi.Mode().String(),
		LineCount: lines,
		WordCount: words,
	}, nil
}

// Analyze analyzes the given text content and returns detailed TextInfo.
func Analyze(content string, filePath string /* optional, for context */) (*ir.TextInfo, error) {
	textInfo := &ir.TextInfo{
		FileInfo: ir.FileInfo{ // Initialize embedded FileInfo
			// Path and Name would ideally be set if filePath is known and relevant
			// For pure text analysis, these might be blank or derived.
			// Content is not part of FileInfo in this context, TextInfo handles it.
			LineCount: len(strings.Split(content, "\n")),
			WordCount: len(strings.Fields(content)),
		},
		CharCount: len(content),
		// Keywords will be populated below
		Readability: &ir.ReadabilityScores{},
		Sentiment:   &ir.SentimentScores{},
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
		kwList = append(kwList, ir.KeywordFrequency{Keyword: k, Frequency: v})
	}
	sort.Slice(kwList, func(i, j int) bool {
		return kwList[i].Frequency > kwList[j].Frequency
	})

	// Store top N keywords (e.g., top 10)
	topN := 10
	if len(kwList) < topN {
		topN = len(kwList)
	}
	textInfo.Keywords = kwList[:topN]

	// TODO: Implement Readability Scores (e.g., Flesch-Kincaid, Gunning Fog)
	// This would typically involve external libraries or complex calculations.
	// textInfo.Readability.FleschKincaidGradeLevel = ...

	// TODO: Implement Sentiment Analysis
	// This would also typically involve external libraries.
	// textInfo.Sentiment.Polarity = ...
	// textInfo.Sentiment.Subjectivity = ...

	if filePath != "" {
		textInfo.FileInfo.Path = filePath
		// Potentially set FileInfo.Name as well
	}

	return textInfo, nil
}
