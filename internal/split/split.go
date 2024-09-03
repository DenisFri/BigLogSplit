package split

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"BigLogSplit/internal/config"
)

func SplitFile(cfg config.Config, updateProgress func(float64)) error {
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

	partNumber := 1
	buffer := make([]byte, 1024*1024) // 1 MB buffer

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

			offset, err := file.Seek(0, io.SeekCurrent)
			if err != nil {
				partFile.Close()
				return fmt.Errorf("error getting file offset: %w", err)
			}
			progressBar := float64(offset) / float64(totalSize)

			time.Sleep(10 * time.Millisecond) // Slow down the progress bar for visual effect

			updateProgress(progressBar)
		}
		writer.Flush()
		partFile.Close()

		if partSize == 0 || partSize < int64(cfg.MaxSizeMB*1024*1024) {
			break
		}

		partNumber++
	}

	updateProgress(1.0) // Mark progress as 100% complete
	return nil
}
