package render

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Jangmyun/gh-contrib-graph/internal/contrib"
	"github.com/Jangmyun/gh-contrib-graph/internal/github"
)

const paddingX = 2

// boxOverhead is everything around the grid's own cells in the widest line:
// 2 border chars + 4 padding chars (paddingX*2) + the row-label column.
const boxOverhead = 2 + paddingX*2 + rowLabelWidth

// Render composes the full bordered box: border-embedded title, contribution
// stats, the heatmap grid, and the color-level legend, sized to fit within
// the current terminal width.
func Render(login string, cal *github.ContributionCalendar, longest, current int) string {
	return render(login, cal, longest, current, TerminalWidth())
}

func render(login string, cal *github.ContributionCalendar, longest, current, width int) string {
	title := dimStyle.Render("GitHub Contributions for ") + usernameStyle.Render(login)

	stats := lipgloss.JoinVertical(lipgloss.Left,
		numberStyle.Render(strconv.Itoa(cal.TotalContributions))+" contributions in the last year",
		"Longest Streak: "+streakStyle.Render(fmt.Sprintf("%d days", longest))+" \U0001F3C6",
		"Current Streak: "+streakStyle.Render(fmt.Sprintf("%d days", current))+" \U0001F525",
	)

	var days []github.ContributionDay
	for _, w := range cal.Weeks {
		days = append(days, w.ContributionDays...)
	}
	max := contrib.MaxCount(days)

	visibleWeeks := fitWeeks(cal.Weeks, width)

	content := lipgloss.JoinVertical(lipgloss.Center, stats, "", Grid(visibleWeeks, max), "", legendLine())

	return box(title, content)
}

// fitWeeks trims weeks (each 2 display cells wide) down to however many of
// the most recent ones fit within width, so the box never exceeds the
// terminal's columns. Older weeks are dropped first; the most recent week
// is always kept.
func fitWeeks(weeks []github.ContributionWeek, width int) []github.ContributionWeek {
	maxCells := (width - boxOverhead) / 2
	if maxCells < 1 {
		maxCells = 1
	}
	if maxCells >= len(weeks) {
		return weeks
	}
	return weeks[len(weeks)-maxCells:]
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
