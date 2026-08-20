package render

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/Jangmyun/gh-contrib-graph/internal/github"
)

var ansiRE = regexp.MustCompile("\x1b\\[[0-9;]*m")

func stripANSI(s string) string {
	return ansiRE.ReplaceAllString(s, "")
}

func fixtureCalendar() *github.ContributionCalendar {
	mkWeek := func(startDate string, counts [7]int) github.ContributionWeek {
		days := make([]github.ContributionDay, 7)
		for wd := 0; wd < 7; wd++ {
			days[wd] = github.ContributionDay{
				Date:              startDate,
				ContributionCount: counts[wd],
				Weekday:           wd,
			}
		}
		return github.ContributionWeek{ContributionDays: days}
	}
	// Enough weeks that the grid (not the title or stat lines) is the widest
	// block, matching real-world rendering where ~52 weeks always dominate.
	weeks := make([]github.ContributionWeek, 0, 20)
	dates := []string{
		"2026-01-04", "2026-01-11", "2026-01-18", "2026-01-25",
		"2026-02-01", "2026-02-08", "2026-02-15", "2026-02-22",
		"2026-03-01", "2026-03-08", "2026-03-15", "2026-03-22",
		"2026-04-01", "2026-04-08", "2026-04-15", "2026-04-22",
		"2026-05-01", "2026-05-08", "2026-05-15", "2026-05-22",
	}
	for i, date := range dates {
		weeks = append(weeks, mkWeek(date, [7]int{i % 3, (i + 1) % 4, 0, (i + 2) % 5, 0, i % 2, 0}))
	}

	return &github.ContributionCalendar{
		TotalContributions: 42,
		Weeks:              weeks,
	}
}

func TestRenderStructure(t *testing.T) {
	// Fixed, generous width so the grid isn't truncated by the test runner's
	// own (possibly absent) terminal size — width-fitting is covered by
	// TestFitWeeks / TestRenderFitsWidth below.
	out := render("testuser", fixtureCalendar(), 5, 2, 200)
	plain := stripANSI(out)
	lines := strings.Split(plain, "\n")

	if !strings.Contains(plain, "GitHub Contributions for testuser") {
		t.Errorf("missing title, got:\n%s", plain)
	}
	if !strings.Contains(plain, "42 contributions in the last year") {
		t.Errorf("missing total contributions line")
	}
	if !strings.Contains(plain, "Longest Streak: 5 days") {
		t.Errorf("missing longest streak line")
	}
	if !strings.Contains(plain, "Current Streak: 2 days") {
		t.Errorf("missing current streak line")
	}
	if !strings.Contains(plain, "Less") || !strings.Contains(plain, "More") {
		t.Errorf("missing legend")
	}
	if !strings.HasPrefix(lines[0], "╭") || !strings.HasSuffix(lines[0], "╮") {
		t.Errorf("top border malformed: %q", lines[0])
	}
	last := lines[len(lines)-1]
	if !strings.HasPrefix(last, "╰") || !strings.HasSuffix(last, "╯") {
		t.Errorf("bottom border malformed: %q", last)
	}

	// Every line (a fully-boxed rectangle) must have the same display width
	// (using display-cell width, since emoji render as double-width).
	width := lipgloss.Width(lines[0])
	for i, l := range lines {
		if got := lipgloss.Width(l); got != width {
			t.Errorf("line %d width = %d, want %d: %q", i, got, width, l)
		}
	}

	// Weekday row labels: only Mon/Wed/Fri should appear as a standalone
	// label, each exactly once; the other four weekdays stay unlabeled.
	labelRE := regexp.MustCompile(`\b(Mon|Tue|Wed|Thu|Fri|Sat|Sun)\b`)
	counts := map[string]int{}
	for _, l := range lines {
		for _, m := range labelRE.FindAllString(l, -1) {
			counts[m]++
		}
	}
	for _, want := range []string{"Mon", "Wed", "Fri"} {
		if counts[want] != 1 {
			t.Errorf("expected exactly one %s row label, got %d", want, counts[want])
		}
	}
	for _, unwanted := range []string{"Tue", "Thu", "Sat", "Sun"} {
		if counts[unwanted] != 0 {
			t.Errorf("unexpected %s row label found (%d occurrences)", unwanted, counts[unwanted])
		}
	}
}

func TestFitWeeks(t *testing.T) {
	weeks := fixtureCalendar().Weeks // 20 weeks

	cases := []struct {
		name  string
		width int
		want  int
	}{
		{"plenty of room keeps every week", 300, len(weeks)},
		{"narrow terminal keeps only the most recent weeks", 40, (40 - boxOverhead) / 2},
		{"absurdly narrow terminal keeps at least one week", 5, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := fitWeeks(weeks, c.width)
			if len(got) != c.want {
				t.Errorf("fitWeeks width=%d returned %d weeks, want %d", c.width, len(got), c.want)
			}
			if len(got) > 0 && len(weeks) > 0 && got[len(got)-1].ContributionDays[0].Date != weeks[len(weeks)-1].ContributionDays[0].Date {
				t.Errorf("fitWeeks must keep the most recent week last")
			}
		})
	}
}

// TestRenderFitsWidth is the regression test for the reported bug: at a
// realistic (narrow) terminal width, every rendered line — including the
// border — must stay within that width, so the terminal never soft-wraps
// mid-box and mangles the border.
func TestRenderFitsWidth(t *testing.T) {
	for _, width := range []int{60, 80, 120, 200} {
		t.Run(fmt.Sprintf("width=%d", width), func(t *testing.T) {
			out := render("testuser", fixtureCalendar(), 5, 2, width)
			plain := stripANSI(out)
			for i, l := range strings.Split(plain, "\n") {
				if got := lipgloss.Width(l); got > width {
					t.Errorf("line %d is %d cells wide, exceeds terminal width %d: %q", i, got, width, l)
				}
			}
		})
	}
}
