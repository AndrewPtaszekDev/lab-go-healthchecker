package internal

import (
	"fmt"
	"slices"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/bubbles/spinner"
)

const autoRefreshRate = time.Second * 25
const refreshDuration = time.Second

type model struct {
	services   []service
	ch         <-chan service
	cursor     int
	refreshing bool
	spinner    spinner.Model
}

type tickMsg time.Time
type refreshStartedMsg struct{ ch <-chan service }
type refreshDoneMsg bool
type checkResultMsg service
type allChecksDoneMsg bool

func (m *model) apply(s service) {
	m.services[s.id] = s
}

// trigger refresh every 'autorefreshRate' seconds
func doTick() tea.Cmd {
	return tea.Tick(autoRefreshRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func alphaSort(services []service) []service {
	slices.SortFunc(services, func(a, b service) int {
		return strings.Compare(a.name, b.name)
	})
	return services
}

func checkForUpdates(ch <-chan service) tea.Cmd {

	return func() tea.Msg {
		s, ok := <-ch // block for exactly one result
		if !ok {
			return allChecksDoneMsg(true) // channel closed = work finished
		}
		return checkResultMsg(s)
	}
}

func refreshServices(services []service) tea.Cmd {
	return func() tea.Msg {
		return refreshStartedMsg{ch: processHeathchecks(services)}
	}
}

// formatElapsed renders a duration at second granularity (no ms/µs),
// trimming a trailing "0s" so values read like "3m35s" or "4h19m".
func formatElapsed(d time.Duration) string {
	s := d.Round(time.Second).String()
	if strings.HasSuffix(s, "m0s") {
		s = strings.TrimSuffix(s, "0s")
	}
	return s
}

func stateStyle(state healthState) lipgloss.Style {
	switch state {
	case Healthy:
		return healthyStyle
	case Stale:
		return staleStyle
	case Unhealthy:
		return unhealthyStyle
	default:
		return unknownStyle
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "space", "r":
			if m.refreshing {
				return m, nil
			}
			return m, refreshServices(m.services)

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.services)-1 {
				m.cursor++
			}
		}

	case checkResultMsg:

		m.apply(service(msg))
		return m, checkForUpdates(m.ch)
	// updated services channel (m.ch) is drained
	case allChecksDoneMsg:
		return m, nil
	case spinner.TickMsg:
		newSpinner, spinnerCmd := m.spinner.Update(msg)
		m.spinner = newSpinner
		if spinnerCmd == nil {
			return m, nil
		}
		return m, func() tea.Msg { return spinnerCmd() }

	case refreshStartedMsg:
		m.ch = msg.ch
		m.refreshing = true
		return m, tea.Batch(
			checkForUpdates(m.ch),
			tea.Tick(refreshDuration, func(t time.Time) tea.Msg { return refreshDoneMsg(true) }),
			func() tea.Msg { return m.spinner.Tick() },
		)

	case refreshDoneMsg:
		m.refreshing = false

	case tickMsg:
		if m.refreshing {
			return m, nil
		}
		return m, tea.Batch(
			doTick(),
			refreshServices(m.services),
		)
	}

	return m, nil
}

func (m model) View() tea.View {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Lab Health"))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Auto-refresh every 25s"))
	b.WriteString("\n\n")

	// Fixed name-column width so every state label lines up vertically.
	nameColWidth := 0
	for _, svc := range m.services {
		// width = cursor icon + space (2) + name + left/right padding (2)
		if w := lipgloss.Width(svc.name) + 4; w > nameColWidth {
			nameColWidth = w
		}
	}

	for i, svc := range m.services {
		cursorIcon := " "
		if i == m.cursor {
			cursorIcon = ">"
		}

		name := fmt.Sprintf("%s %s", cursorIcon, svc.name)
		if i == m.cursor {
			name = selectedStyle.Width(nameColWidth).Render(name)
		} else {
			name = normalStyle.Width(nameColWidth).Render(name)
		}

		elapsedString := formatElapsed(svc.atStateElasped)
		if svc.atStateElasped <= time.Second*10 {
			elapsedString = "moments ago"
		}

		stateText := fmt.Sprintf("%s (%s)", svc.state.String(), elapsedString)
		stateLabel := stateStyle(svc.state).Render(stateText)

		b.WriteString(name)
		b.WriteString("  ")
		b.WriteString(stateLabel)
		b.WriteString("\n")
	}

	if m.refreshing {
		b.WriteString("\n")
		b.WriteString(spinnerStyle.Render(m.spinner.View()))
	} else {
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("space: refresh | q: quit"))
	}

	return tea.NewView(borderStyle.Render(b.String()))
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		refreshServices(m.services),
		doTick(),
	)
}

func initialModel(services []service) tea.Model {
	sorted := alphaSort(services)
	for i := range sorted {
		sorted[i].id = i
	}
	s := spinner.New()
	s.Spinner = spinner.Dot

	return model{
		services:   sorted,
		cursor:     0,
		refreshing: true,
		spinner:    s,
	}
}
