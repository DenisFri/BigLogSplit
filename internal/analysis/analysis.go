package analysis

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type PatternMatch struct {
	Pattern     string
	Count       int
	Lines       []string
	LineNumbers []int
}

type AnalysisResult struct {
	ErrorMatches   map[string]*PatternMatch
	CustomMatches  map[string]*PatternMatch
	TotalLines     int
	ProcessedBytes int64
}

func NewAnalysisResult() *AnalysisResult {
	return &AnalysisResult{
		ErrorMatches:  make(map[string]*PatternMatch),
		CustomMatches: make(map[string]*PatternMatch),
	}
}

func AnalyzeFile(filePath string, errorPatterns, customPatterns []*regexp.Regexp, maxMatchesPerPattern int) (*AnalysisResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file for analysis: %w", err)
	}
	defer file.Close()

	result := NewAnalysisResult()

	// Initialize pattern matches
	for _, pattern := range errorPatterns {
		result.ErrorMatches[pattern.String()] = &PatternMatch{
			Pattern:     pattern.String(),
			Lines:       make([]string, 0, maxMatchesPerPattern),
			LineNumbers: make([]int, 0, maxMatchesPerPattern),
		}
	}

	for _, pattern := range customPatterns {
		result.CustomMatches[pattern.String()] = &PatternMatch{
			Pattern:     pattern.String(),
			Lines:       make([]string, 0, maxMatchesPerPattern),
			LineNumbers: make([]int, 0, maxMatchesPerPattern),
		}
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		result.TotalLines++
		result.ProcessedBytes += int64(len(line) + 1) // +1 for newline

		// Check for error patterns
		for _, pattern := range errorPatterns {
			if pattern.MatchString(line) {
				match := result.ErrorMatches[pattern.String()]
				match.Count++

				// Store the line if under max matches limit
				if len(match.Lines) < maxMatchesPerPattern {
					match.Lines = append(match.Lines, line)
					match.LineNumbers = append(match.LineNumbers, lineNum)
				}
			}
		}

		// Check for custom patterns
		for _, pattern := range customPatterns {
			if pattern.MatchString(line) {
				match := result.CustomMatches[pattern.String()]
				match.Count++

				// Store the line if under max matches limit
				if len(match.Lines) < maxMatchesPerPattern {
					match.Lines = append(match.Lines, line)
					match.LineNumbers = append(match.LineNumbers, lineNum)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return result, nil
}

func WriteAnalysisReport(result *AnalysisResult, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating analysis report file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write header
	writer.WriteString("# Log Analysis Report\n\n")
	writer.WriteString(fmt.Sprintf("Total lines processed: %d\n", result.TotalLines))
	writer.WriteString(fmt.Sprintf("Total bytes processed: %s\n\n", formatSize(result.ProcessedBytes)))

	// Write error pattern matches
	writer.WriteString("## Error Pattern Matches\n\n")
	if len(result.ErrorMatches) == 0 {
		writer.WriteString("No error patterns defined or matched.\n\n")
	} else {
		for _, match := range result.ErrorMatches {
			writer.WriteString(fmt.Sprintf("### Pattern: `%s`\n", match.Pattern))
			writer.WriteString(fmt.Sprintf("- Total matches: %d\n\n", match.Count))

			if len(match.Lines) > 0 {
				writer.WriteString("#### Sample Matches:\n")
				for i, line := range match.Lines {
					writer.WriteString(fmt.Sprintf("Line %d: `%s`\n", match.LineNumbers[i], truncateString(line, 100)))
				}
				writer.WriteString("\n")
			}
		}
	}

	// Write custom pattern matches
	writer.WriteString("## Custom Pattern Matches\n\n")
	if len(result.CustomMatches) == 0 {
		writer.WriteString("No custom patterns defined or matched.\n\n")
	} else {
		for _, match := range result.CustomMatches {
			writer.WriteString(fmt.Sprintf("### Pattern: `%s`\n", match.Pattern))
			writer.WriteString(fmt.Sprintf("- Total matches: %d\n\n", match.Count))

			if len(match.Lines) > 0 {
				writer.WriteString("#### Sample Matches:\n")
				for i, line := range match.Lines {
					writer.WriteString(fmt.Sprintf("Line %d: `%s`\n", match.LineNumbers[i], truncateString(line, 100)))
				}
				writer.WriteString("\n")
			}
		}
	}

	return writer.Flush()
}

// FilterLine checks if a line should be included based on filter mode and patterns
func FilterLine(line string, mode string, patterns []*regexp.Regexp) bool {
	if len(patterns) == 0 || mode == "none" {
		return true // No filtering
	}

	matches := false
	for _, pattern := range patterns {
		if pattern.MatchString(line) {
			matches = true
			break
		}
	}

	if mode == "include" {
		return matches // Include only lines matching patterns
	} else if mode == "exclude" {
		return !matches // Exclude lines matching patterns
	}

	return true // Default: include all
}

// Helper function to truncate long strings
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Helper to format file size
func formatSize(sizeBytes int64) string {
	const (
		_          = iota
		KB float64 = 1 << (10 * iota)
		MB
		GB
	)

	var size float64 = float64(sizeBytes)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", size/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", size/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", size/KB)
	default:
		return fmt.Sprintf("%d bytes", sizeBytes)
	}
}
