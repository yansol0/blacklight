package tui

import (
	"fmt"
	"time"

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
	Phase        string // Unauth, Auth, Bypass <Header>, Done
}

type doneMsg struct{}

type Model struct {
	spinner   spinner.Model
	progress  progress.Model
	status    string
	detail    string
	percent   float64
	updatesCh <-chan ProgressUpdate
	doneCh    <-chan struct{}
	quitting  bool
}

func NewModel(updates <-chan ProgressUpdate, done <-chan struct{}) Model {
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

func waitForDone(ch <-chan struct{}) tea.Cmd {
	return func() tea.Msg {
		<-ch
		return doneMsg{}
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
	case doneMsg:
		m.status = "Done"
		m.detail = "All endpoints tested"
		m.percent = 1.0
		m.quitting = true
		return m, tea.Tick(500*time.Millisecond, func(time.Time) tea.Msg { return tea.Quit })
	}
	return m, nil
}

func (m Model) View() string {
	bar := m.progress.ViewAs(m.percent)
	return fmt.Sprintf("%s %s\n%s\n%s\n(press q to quit)", m.spinner.View(), m.status, bar, m.detail)
}

func Run(updates <-chan ProgressUpdate, done <-chan struct{}) error {
	p := tea.NewProgram(NewModel(updates, done))
	_, err := p.Run()
	return err
}
