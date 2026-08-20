package render

import "github.com/charmbracelet/lipgloss"

// levelColors mirrors GitHub's own dark-theme contribution calendar palette.
var levelColors = [5]lipgloss.Color{
	lipgloss.Color("#161b22"),
	lipgloss.Color("#0e4429"),
	lipgloss.Color("#006d32"),
	lipgloss.Color("#26a641"),
	lipgloss.Color("#39d353"),
}

var (
	usernameStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff"))
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#8b949e"))
	numberStyle   = lipgloss.NewStyle().Bold(true)
	streakStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#39c5cf"))
	borderStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#8b949e"))
)

func cellStyle(level int) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(levelColors[level])
}
