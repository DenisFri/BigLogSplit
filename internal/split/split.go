package split

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"BigLogSplit/internal/analysis"
	"BigLogSplit/internal/config"
	"BigLogSplit/internal/ui"
)

func SplitFile(cfg config.RuntimeConfig, updateProgress func(interface{})) error {
	file, err := os.Open(cfg.FilePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	totalSize := fileInfo.Size()

	// Ensure the output directory exists
	if err := os.MkdirAll(cfg.OutputFolder, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Initialize analysis if enabled
	var analysisResult *analysis.AnalysisResult
	if cfg.Analysis.Enabled {
		analysisResult = analysis.NewAnalysisResult()
	}

	partNumber := 1
	totalProcessed := int64(0)
	reader := bufio.NewReader(file)

	for {
		partFilePath := filepath.Join(cfg.OutputFolder, fmt.Sprintf("part%d.log", partNumber))
		partFile, err := os.Create(partFilePath)
		if err != nil {
			return fmt.Errorf("error creating part file: %w", err)
		}

		writer := bufio.NewWriter(partFile)
		partSize := int64(0)
		lineNumber := 0

		for partSize < int64(cfg.MaxSizeMB*1024*1024) {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				partFile.Close()
				return fmt.Errorf("error reading file: %w", err)
			}

			if line == "" && err == io.EOF {
				break
			}

			lineNumber++
			lineSize := int64(len(line))

			// Apply filtering if enabled
			if cfg.Filtering.Mode != config.FilterModeNone && len(cfg.FilterPatternRegexps) > 0 {
				shouldInclude := analysis.FilterLine(
					line,
					string(cfg.Filtering.Mode),
					cfg.FilterPatternRegexps,
				)

				if !shouldInclude {
					// Skip this line, but still count it as processed
					totalProcessed += lineSize
					continue
				}
			}

			// Write the line to the output file
			if _, err := writer.WriteString(line); err != nil {
				partFile.Close()
				return fmt.Errorf("error writing to part file: %w", err)
			}
			writer.Flush()

			partSize += lineSize
			totalProcessed += lineSize

			// If analysis is enabled, check patterns
			if cfg.Analysis.Enabled {
				analysisResult.TotalLines++
				analysisResult.ProcessedBytes += lineSize

				// Check error patterns
				for _, pattern := range cfg.ErrorPatternRegexps {
					if pattern.MatchString(line) {
						match, exists := analysisResult.ErrorMatches[pattern.String()]
						if !exists {
							match = &analysis.PatternMatch{
								Pattern:     pattern.String(),
								Lines:       make([]string, 0, 10),
								LineNumbers: make([]int, 0, 10),
							}
							analysisResult.ErrorMatches[pattern.String()] = match
						}

						match.Count++
						if len(match.Lines) < 10 { // Store up to 10 sample lines
							match.Lines = append(match.Lines, line)
							match.LineNumbers = append(match.LineNumbers, lineNumber)
						}
					}
				}

				// Check custom patterns
				for _, pattern := range cfg.CustomPatternRegexps {
					if pattern.MatchString(line) {
						match, exists := analysisResult.CustomMatches[pattern.String()]
						if !exists {
							match = &analysis.PatternMatch{
								Pattern:     pattern.String(),
								Lines:       make([]string, 0, 10),
								LineNumbers: make([]int, 0, 10),
							}
							analysisResult.CustomMatches[pattern.String()] = match
						}

						match.Count++
						if len(match.Lines) < 10 { // Store up to 10 sample lines
							match.Lines = append(match.Lines, line)
							match.LineNumbers = append(match.LineNumbers, lineNumber)
						}
					}
				}
			}

			progressBar := float64(totalProcessed) / float64(totalSize)

			// Update progress less frequently for better performance
			if totalProcessed%int64(1024*1024) == 0 { // Update every ~1MB
				time.Sleep(10 * time.Millisecond) // Small delay for UI updates

				// Send both progress percentage and processed bytes
				updateProgress(ui.ProgressUpdate{
					Percent:        progressBar,
					ProcessedBytes: totalProcessed,
				})
			}

			if err == io.EOF {
				break
			}
		}

		writer.Flush()
		partFile.Close()

		if partSize == 0 {
			// Remove empty part file
			os.Remove(partFilePath)
			break
		}

		if partSize < int64(cfg.MaxSizeMB*1024*1024) {
			break
		}

		partNumber++
	}

	// Write analysis report if enabled
	if cfg.Analysis.Enabled && cfg.Analysis.OutputFile != "" {
		outputPath := cfg.Analysis.OutputFile
		if !filepath.IsAbs(outputPath) {
			outputPath = filepath.Join(cfg.OutputFolder, outputPath)
		}

		if err := analysis.WriteAnalysisReport(analysisResult, outputPath); err != nil {
			return fmt.Errorf("error writing analysis report: %w", err)
		}
	}

	// Mark progress as 100% complete with total size
	updateProgress(ui.ProgressUpdate{
		Percent:        1.0,
		ProcessedBytes: totalSize,
	})
	return nil
}
