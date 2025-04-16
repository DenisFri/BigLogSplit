package main

import (
	"fmt"
	"os"
	"path/filepath"

	"BigLogSplit/internal/config"
	"BigLogSplit/internal/split"
	"BigLogSplit/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg, err := config.ReadConfig("config.json")
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	// Get the file size for UI initialization
	fileInfo, err := os.Stat(cfg.FilePath)
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}
	totalSize := fileInfo.Size()

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(cfg.OutputFolder, 0755); err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	// Set default analysis output file if enabled but not specified
	if cfg.Analysis.Enabled && cfg.Analysis.OutputFile == "" {
		cfg.Analysis.OutputFile = filepath.Join(cfg.OutputFolder, "analysis_report.md")
	}

	// Initialize the UI with the total file size
	p := ui.NewModel(totalSize)
	prog := tea.NewProgram(p)

	go func() {
		if err := split.SplitFile(cfg, func(update interface{}) {
			prog.Send(update)
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
