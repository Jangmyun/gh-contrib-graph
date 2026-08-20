package render

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Jangmyun/gh-contrib-graph/internal/github"
)

// Render composes the full bordered box: border-embedded title, contribution
// stats, the heatmap grid, and the color-level legend.
func Render(login string, cal *github.ContributionCalendar, longest, current int) string {
	title := dimStyle.Render("GitHub Contributions for ") + usernameStyle.Render(login)

	stats := lipgloss.JoinVertical(lipgloss.Left,
		numberStyle.Render(strconv.Itoa(cal.TotalContributions))+" contributions in the last year",
		"Longest Streak: "+streakStyle.Render(fmt.Sprintf("%d days", longest))+" \U0001F3C6",
		"Current Streak: "+streakStyle.Render(fmt.Sprintf("%d days", current))+" \U0001F525",
	)

	content := lipgloss.JoinVertical(lipgloss.Center, stats, "", Grid(cal), "", legendLine())

	return box(title, content)
}

func legendLine() string {
	var b strings.Builder
	b.WriteString(dimStyle.Render("Less "))
	for lvl := 0; lvl < 5; lvl++ {
		b.WriteString(cellStyle(lvl).Render(cellGlyph))
		b.WriteString(" ")
	}
	b.WriteString(dimStyle.Render("More"))
	return b.String()
}

// box wraps content in a rounded border with title spliced into the top edge.
func box(title, content string) string {
	const paddingX = 2

	contentWidth := lipgloss.Width(content)
	titleWidth := lipgloss.Width(title)

	innerWidth := contentWidth + paddingX*2
	if minWidth := titleWidth + 4; innerWidth < minWidth {
		innerWidth = minWidth
	}

	dashesTotal := innerWidth - titleWidth - 2
	leftDashes := dashesTotal / 2
	rightDashes := dashesTotal - leftDashes

	top := borderStyle.Render("╭"+strings.Repeat("─", leftDashes)) +
		" " + title + " " +
		borderStyle.Render(strings.Repeat("─", rightDashes)+"╮")
	bottom := borderStyle.Render("╰" + strings.Repeat("─", innerWidth) + "╯")
	blankRow := borderStyle.Render("│" + strings.Repeat(" ", innerWidth) + "│")

	body := lipgloss.NewStyle().Width(contentWidth).Align(lipgloss.Center).Render(content)

	var b strings.Builder
	b.WriteString(top + "\n")
	b.WriteString(blankRow + "\n")
	for _, line := range strings.Split(body, "\n") {
		b.WriteString(borderStyle.Render("│") + strings.Repeat(" ", paddingX) + line + strings.Repeat(" ", paddingX) + borderStyle.Render("│") + "\n")
	}
	b.WriteString(blankRow + "\n")
	b.WriteString(bottom)
	return b.String()
}
