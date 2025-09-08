package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// ProgressUpdate is sent by the tester to update the TUI
type ProgressUpdate struct {
	CurrentIndex int
	Total        int
	Method       string
	URL          string
	Phase        string // Unauth, Auth, Bypass <Header>
}

// Summary is sent when all work is finished
type Summary struct {
	BypassHits     int
	IDORCandidates int
}

type Model struct {
	spinner   spinner.Model
	progress  progress.Model
	status    string
	detail    string
	percent   float64
	updatesCh <-chan ProgressUpdate
	doneCh    <-chan Summary
	quitting  bool
	showHint  bool
}

func NewModel(updates <-chan ProgressUpdate, done <-chan Summary) Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	pr := progress.New(progress.WithScaledGradient("#3B82F6", "#10B981"))
	return Model{
		spinner:   sp,
		progress:  pr,
		status:    "Starting...",
		detail:    "",
		percent:   0,
		updatesCh: updates,
		doneCh:    done,
		showHint:  true,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, waitForUpdate(m.updatesCh), waitForDone(m.doneCh))
}

func waitForUpdate(ch <-chan ProgressUpdate) tea.Cmd {
	return func() tea.Msg {
		upd, ok := <-ch
		if !ok {
			return nil
		}
		return upd
	}
}

func waitForDone(ch <-chan Summary) tea.Cmd {
	return func() tea.Msg {
		s := <-ch
		return s
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case ProgressUpdate:
		m.status = fmt.Sprintf("Testing %d/%d", msg.CurrentIndex, msg.Total)
		m.detail = fmt.Sprintf("%s %s [%s]", msg.Method, msg.URL, msg.Phase)
		if msg.Total > 0 {
			m.percent = float64(msg.CurrentIndex) / float64(msg.Total)
		}
		return m, tea.Batch(waitForUpdate(m.updatesCh), m.spinner.Tick)
	case Summary:
		m.status = "Done"
		m.detail = fmt.Sprintf("Auth bypasses: %d | IDOR candidates: %d", msg.BypassHits, msg.IDORCandidates)
		m.percent = 1.0
		m.showHint = false
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) View() string {
	bar := m.progress.ViewAs(m.percent)
	hint := ""
	if m.showHint {
		hint = "\n(press q to quit)"
	}
	return fmt.Sprintf("%s %s\n%s\n%s%s", m.spinner.View(), m.status, bar, m.detail, hint)
}

func Run(updates <-chan ProgressUpdate, done <-chan Summary) error {
	p := tea.NewProgram(NewModel(updates, done))
	_, err := p.Run()
	return err
}
