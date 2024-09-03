package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	progress progress.Model
	percent  float64
	done     bool
	message  string
}

func NewModel() Model {
	p := progress.NewModel()
	return Model{progress: p, percent: 0, done: false}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.done {
			return m, tea.Quit
		}
	case float64:
		m.percent = msg
		if m.percent >= 1.0 {
			m.done = true
			m.message = "File splitting complete! Press any key to exit."
			return m, nil
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.done {
		return fmt.Sprintf("\n%s\n", m.message)
	}
	return fmt.Sprintf("\nProgress: %s\n", m.progress.ViewAs(m.percent))
}
