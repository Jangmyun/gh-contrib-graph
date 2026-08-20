package contrib

import (
	"testing"

	"github.com/Jangmyun/gh-contrib-graph/internal/github"
)

// days builds a slice of ContributionDay from raw counts, in chronological order.
func days(counts ...int) []github.ContributionDay {
	out := make([]github.ContributionDay, len(counts))
	for i, c := range counts {
		out[i] = github.ContributionDay{ContributionCount: c}
	}
	return out
}

func TestLongestStreak(t *testing.T) {
	cases := []struct {
		name string
		in   []github.ContributionDay
		want int
	}{
		{"empty", nil, 0},
		{"all zero", days(0, 0, 0), 0},
		{"all nonzero", days(1, 2, 3), 3},
		{"broken mid-way", days(1, 1, 0, 1, 1, 1, 0, 1), 3},
		{"trailing zero today", days(1, 1, 1, 0), 3},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := LongestStreak(c.in); got != c.want {
				t.Errorf("LongestStreak = %d, want %d", got, c.want)
			}
		})
	}
}

func TestCurrentStreak(t *testing.T) {
	cases := []struct {
		name string
		in   []github.ContributionDay
		want int
	}{
		{"empty", nil, 0},
		{"all zero", days(0, 0, 0), 0},
		{"all nonzero", days(1, 2, 3), 3},
		{"trailing zero today still counts prior streak", days(1, 1, 1, 0), 3},
		{"broken then resumed", days(1, 1, 0, 1, 1, 1), 3},
		{"broken then zero today", days(1, 1, 1, 0, 0), 0},
		{"single day with contribution", days(1), 1},
		{"single day no contribution", days(0), 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := CurrentStreak(c.in); got != c.want {
				t.Errorf("CurrentStreak = %d, want %d", got, c.want)
			}
		})
	}
}
