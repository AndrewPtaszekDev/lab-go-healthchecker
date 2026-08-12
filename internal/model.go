package internal

import (
	"time"
	"slices"
	"strings"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	tea "charm.land/bubbletea/v2"
)

const autoRefreshRate = time.Second * 25

type model struct {
	services   []service
	cursor     int
	refreshing bool
}

type tickMsg time.Time
type refreshDoneMsg []service

type palette struct {
	title     string
	selected  string
	normal    string
	muted     string
	border    string
	help      string
	healthy   string
	stale     string
	unhealthy string
	unknown   string
}

var theme = palette{
	title:     "#f4efe6",
	selected:  "#c4934f",
	normal:    "#d9cfc6",
	muted:     "#6e5a61",
	border:    "#463139",
	help:      "#ce8aa3",
	healthy:   "#7dbb7f",
	stale:     "#c4934f",
	unhealthy: "#c0685b",
	unknown:   "#6e5a61",
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(theme.title))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.selected)).
			Background(lipgloss.Color("#334155")).
			Padding(0, 1)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.normal)).
			Padding(0, 1)

	mutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.muted))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.help))

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.border)).
			Padding(1, 2)

	healthyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.healthy))
	staleStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.stale))
	unhealthyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.unhealthy))
	unknownStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.unknown))
)

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

func refreshServices(services []service) tea.Cmd {
	return func() tea.Msg {
		processed := processHeathchecks(services)
		sorted := alphaSort(processed)

		return refreshDoneMsg(sorted)
	}
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
		case "space", "enter", "r":
			if m.refreshing {
				return m, nil
			}

			m.refreshing = true
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

	case refreshDoneMsg:
		m.refreshing = false
		m.services = []service(msg)

	case tickMsg:
		if m.refreshing {
			return m, nil
		}
		m.refreshing = true
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

	for i, svc := range m.services {
		cursorIcon := " "
		if i == m.cursor {
			cursorIcon = ">"
		}

		name := fmt.Sprintf("%s %s", cursorIcon, svc.name)
		if i == m.cursor {
			name = selectedStyle.Render(name)
		} else {
			name = normalStyle.Render(name)
		}

		elapsedString := svc.atStateElasped.String()
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
		b.WriteString(helpStyle.Render("refreshing..."))
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
	return model{
		services:   services,
		cursor:     0,
		refreshing: true,
	}
}
