package split

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"BigLogSplit/internal/config"
	"BigLogSplit/internal/ui"
)

func SplitFile(cfg config.Config, updateProgress func(interface{})) error {
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

	partNumber := 1
	buffer := make([]byte, 1024*1024) // 1 MB buffer
	totalProcessed := int64(0)

	for {
		partFilePath := filepath.Join(cfg.OutputFolder, fmt.Sprintf("part%d.log", partNumber))
		partFile, err := os.Create(partFilePath)
		if err != nil {
			return fmt.Errorf("error creating part file: %w", err)
		}

		writer := bufio.NewWriter(partFile)
		partSize := int64(0)
		for partSize < int64(cfg.MaxSizeMB*1024*1024) {
			n, err := file.Read(buffer)
			if err != nil && err != io.EOF {
				partFile.Close()
				return fmt.Errorf("error reading file: %w", err)
			}
			if n == 0 {
				break
			}

			if _, err := writer.Write(buffer[:n]); err != nil {
				partFile.Close()
				return fmt.Errorf("error writing to part file: %w", err)
			}
			writer.Flush()

			partSize += int64(n)
			totalProcessed += int64(n)

			progressBar := float64(totalProcessed) / float64(totalSize)

			time.Sleep(10 * time.Millisecond) // Slow down the progress bar for visual effect

			// Send both progress percentage and processed bytes
			updateProgress(ui.ProgressUpdate{
				Percent:        progressBar,
				ProcessedBytes: totalProcessed,
			})
		}
		writer.Flush()
		partFile.Close()

		if partSize == 0 || partSize < int64(cfg.MaxSizeMB*1024*1024) {
			break
		}

		partNumber++
	}

	// Mark progress as 100% complete with total size
	updateProgress(ui.ProgressUpdate{
		Percent:        1.0,
		ProcessedBytes: totalSize,
	})
	return nil
}
