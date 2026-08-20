package render

import (
	"strings"
	"time"

	"github.com/Jangmyun/gh-contrib-graph/internal/contrib"
	"github.com/Jangmyun/gh-contrib-graph/internal/github"
)

const (
	rowLabelWidth = 4
	cellGlyph     = "■"
)

var weekdayLabels = map[int]string{1: "Mon", 3: "Wed", 5: "Fri"}

// Grid renders the month-header row plus the 7 weekday rows of the heatmap.
func Grid(cal *github.ContributionCalendar) string {
	weeks := cal.Weeks

	var days []github.ContributionDay
	for _, w := range weeks {
		days = append(days, w.ContributionDays...)
	}
	max := contrib.MaxCount(days)

	lines := make([]string, 0, 8)
	lines = append(lines, monthHeader(weeks))
	for weekday := 0; weekday < 7; weekday++ {
		lines = append(lines, weekdayRow(weeks, weekday, max))
	}
	return strings.Join(lines, "\n")
}

func monthHeader(weeks []github.ContributionWeek) string {
	buf := []rune(strings.Repeat(" ", len(weeks)*2))
	lastMonth := time.Month(0)
	for i, w := range weeks {
		if len(w.ContributionDays) == 0 {
			continue
		}
		t, err := time.Parse("2006-01-02", w.ContributionDays[0].Date)
		if err != nil {
			continue
		}
		if t.Month() == lastMonth {
			continue
		}
		lastMonth = t.Month()
		for j, r := range t.Format("Jan") {
			pos := i*2 + j
			if pos >= len(buf) {
				break
			}
			buf[pos] = r
		}
	}
	return padRight("", rowLabelWidth) + dimStyle.Render(string(buf))
}

func weekdayRow(weeks []github.ContributionWeek, weekday int, max int) string {
	var b strings.Builder
	b.WriteString(dimStyle.Render(padRight(weekdayLabels[weekday], rowLabelWidth)))
	for _, w := range weeks {
		count := -1
		for _, d := range w.ContributionDays {
			if d.Weekday == weekday {
				count = d.ContributionCount
				break
			}
		}
		if count < 0 {
			b.WriteString("  ")
			continue
		}
		lvl := contrib.Level(count, max)
		b.WriteString(cellStyle(lvl).Render(cellGlyph))
		b.WriteString(" ")
	}
	return b.String()
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
