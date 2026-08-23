package internal

import "charm.land/lipgloss/v2"

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
	spinnerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.help))
)
