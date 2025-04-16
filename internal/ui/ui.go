package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProgressUpdate contains progress information to update the UI
type ProgressUpdate struct {
	Percent        float64
	ProcessedBytes int64
}

type Model struct {
	progress      progress.Model
	spinner       spinner.Model
	percent       float64
	done          bool
	message       string
	startTime     time.Time
	totalSize     int64
	processedSize int64
	useSpinner    bool
}

func NewModel(totalSize int64) Model {
	p := progress.NewModel()
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Use spinner for files smaller than 10MB
	useSpinner := totalSize < 10*1024*1024

	return Model{
		progress:      p,
		spinner:       s,
		percent:       0,
		done:          false,
		startTime:     time.Now(),
		totalSize:     totalSize,
		processedSize: 0,
		useSpinner:    useSpinner,
	}
}

func (m Model) Init() tea.Cmd {
	if m.useSpinner {
		return spinner.Tick
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.done {
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case ProgressUpdate:
		m.percent = msg.Percent
		m.processedSize = msg.ProcessedBytes

		if m.percent >= 1.0 {
			m.done = true
			m.message = "File splitting complete! Press any key to exit."
			return m, nil
		}

		if m.useSpinner {
			return m, spinner.Tick
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.done {
		return fmt.Sprintf("\n%s\n", m.message)
	}

	// Format processed and total size
	processedStr := formatSize(m.processedSize)
	totalStr := formatSize(m.totalSize)

	// Calculate ETA
	eta := calculateETA(m.startTime, m.percent)

	var sb strings.Builder
	sb.WriteString("\n")

	// Show progress info
	if m.useSpinner {
		sb.WriteString(fmt.Sprintf("%s Processing file... %s of %s\n", m.spinner.View(), processedStr, totalStr))
	} else {
		sb.WriteString(fmt.Sprintf("Progress: %s\n", m.progress.ViewAs(m.percent)))
		sb.WriteString(fmt.Sprintf("Processed: %s of %s (%.1f%%)\n", processedStr, totalStr, m.percent*100))
	}

	// Show ETA
	if m.percent > 0 && m.percent < 1 {
		sb.WriteString(fmt.Sprintf("ETA: %s\n", eta))
	}

	return sb.String()
}

// Calculate estimated time remaining
func calculateETA(startTime time.Time, percentComplete float64) string {
	if percentComplete <= 0 {
		return "calculating..."
	}

	elapsed := time.Since(startTime)
	if percentComplete >= 1.0 {
		return "completed"
	}

	// Calculate total expected time based on current progress
	totalExpectedTime := time.Duration(float64(elapsed) / percentComplete)
	timeRemaining := totalExpectedTime - elapsed

	if timeRemaining < time.Minute {
		return fmt.Sprintf("%d seconds", int(timeRemaining.Seconds()))
	} else if timeRemaining < time.Hour {
		return fmt.Sprintf("%d minutes, %d seconds",
			int(timeRemaining.Minutes()),
			int(timeRemaining.Seconds())%60)
	} else {
		return fmt.Sprintf("%d hours, %d minutes",
			int(timeRemaining.Hours()),
			int(timeRemaining.Minutes())%60)
	}
}

// Format file size to human-readable format
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
