package main

import (
	"fmt"
	"os"

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

	p := ui.NewModel()
	prog := tea.NewProgram(p)

	go func() {
		if err := split.SplitFile(cfg, func(percent float64) {
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
