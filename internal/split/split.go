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

// ProcessStatus represents the current operation being performed
type ProcessStatus string

const (
	StatusSplitting ProcessStatus = "Splitting"
	StatusFiltering ProcessStatus = "Filtering"
	StatusAnalyzing ProcessStatus = "Analyzing"
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

	// Determine how many passes we need to make
	totalPasses := 1.0 // Splitting is always required
	if cfg.Filtering.Mode != config.FilterModeNone && len(cfg.FilterPatternRegexps) > 0 {
		totalPasses += 0.2 // Filtering adds 20% to progress calculation
	}
	if cfg.Analysis.Enabled {
		totalPasses += 0.3 // Analysis adds 30% to progress calculation
	}

	// Reset file position
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("error seeking file: %w", err)
	}

	partNumber := 1
	totalProcessed := int64(0)
	reader := bufio.NewReader(file)
	lastUpdateBytes := int64(0)

	// Update UI to show we're starting the splitting process
	updateProgress(ui.StatusUpdate{
		Status: string(StatusSplitting),
	})

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

			// Apply filtering if enabled - now considered part of the splitting process
			shouldInclude := true
			if cfg.Filtering.Mode != config.FilterModeNone && len(cfg.FilterPatternRegexps) > 0 {
				shouldInclude = analysis.FilterLine(
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

			// If analysis is enabled, process it in-line with splitting
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

			// Calculate progress considering all operations
			// Base progress is based on how much file we've processed
			baseProgress := float64(totalProcessed) / float64(totalSize)

			// Scale the base progress to account for the fraction that splitting represents
			scaledProgress := baseProgress / totalPasses

			// Update progress roughly every megabyte processed
			if totalProcessed-lastUpdateBytes >= int64(1024*1024) {
				lastUpdateBytes = totalProcessed

				time.Sleep(10 * time.Millisecond) // Small delay for UI updates

				// Send progress update
				updateProgress(ui.ProgressUpdate{
					Percent:        scaledProgress,
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

	// Splitting is complete, now we're at the base progress level
	splitProgress := 1.0 / totalPasses

	// If analysis report needs to be written, show that as an additional step
	if cfg.Analysis.Enabled && cfg.Analysis.OutputFile != "" {
		// Update status to show we're now analyzing
		updateProgress(ui.StatusUpdate{
			Status: string(StatusAnalyzing),
		})

		// Show progress in the analysis phase
		for i := 0; i < 5; i++ { // Simulate analysis progress steps
			time.Sleep(100 * time.Millisecond)
			analysisProgress := splitProgress + (0.3 * float64(i+1) / 5.0) // Incremental progress in analysis phase

			updateProgress(ui.ProgressUpdate{
				Percent:        analysisProgress,
				ProcessedBytes: totalSize, // Keep the processed bytes at total file size
			})
		}

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
