package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"os"
	"path/filepath"
	"time"
)

type config struct {
	FilePath     string `json:"filePath"`
	MaxSizeMB    int    `json:"maxSizeMB"`
	OutputFolder string `json:"outputFolder"`
}

type model struct {
	progress progress.Model
	percent  float64
	done     bool
	message  string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case float64:
		m.percent = msg
		if m.percent >= 1.0 {
			m.done = true
			m.message = "File splitting complete! Press any key to exit."
			return m, tea.Batch(tea.Printf(m.message), tea.EnterAltScreen)
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.done {
		return fmt.Sprintf("\n%s\n", m.message)
	}
	return fmt.Sprintf("\nProgress: %s\n", m.progress.ViewAs(m.percent))
}

func readConfig(filePath string) (config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return config{}, err
	}
	defer file.Close()

	var cfg config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return config{}, err
	}

	return cfg, nil
}

func splitFile(cfg config, updateProgress func(float64)) error {
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
	buffer := make([]byte, 1024*1024)

	for {
		partFilePath := filepath.Join(cfg.OutputFolder, fmt.Sprintf("part%d.log", partNumber))
		partFile, err := os.Create(partFilePath)
		if err != nil {
			return fmt.Errorf("error creating part file: %w", err)
		}
		defer func(partFile *os.File) {
			err := partFile.Close()
			if err != nil {
				fmt.Println("Error closing part file:", err)
			}
		}(partFile)

		writer := bufio.NewWriter(partFile)
		partSize := int64(0)
		for partSize < int64(cfg.MaxSizeMB*1024*1024) {
			n, err := file.Read(buffer)
			if err != nil && err != io.EOF {
				return fmt.Errorf("error reading file: %w", err)
			}
			if n == 0 {
				break
			}

			writer.Write(buffer[:n])
			writer.Flush() // Flush immediately

			partSize += int64(n)

			offset, err := file.Seek(0, io.SeekCurrent)
			if err != nil {
				return fmt.Errorf("error getting file offset: %w", err)
			}
			progressBar := float64(offset) / float64(totalSize)

			// Add a small delay to slow down the progressBar bar
			time.Sleep(10 * time.Millisecond)

			updateProgress(progressBar)
		}
		writer.Flush()

		if partSize == 0 || partSize < int64(cfg.MaxSizeMB*1024*1024) {
			break
		}

		partNumber++
	}

	updateProgress(1.0) // Mark progress as 100% complete
	return nil
}

func main() {
	cfg, err := readConfig("config.json")
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	p := progress.NewModel()
	m := model{progress: p, percent: 0, done: false}

	prog := tea.NewProgram(m)

	go func() {
		if err := splitFile(cfg, func(percent float64) {
			prog.Send(percent)
		}); err != nil {
			fmt.Println("Error splitting file:", err)
			prog.Send(tea.Quit())
		}
	}()

	if err := prog.Start(); err != nil {
		fmt.Println("Error starting program:", err)
		os.Exit(1)
	}
}
