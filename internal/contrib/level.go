package contrib

import "github.com/Jangmyun/gh-contrib-graph/internal/github"

// MaxCount returns the highest single-day contribution count in days.
func MaxCount(days []github.ContributionDay) int {
	max := 0
	for _, d := range days {
		if d.ContributionCount > max {
			max = d.ContributionCount
		}
	}
	return max
}

// Level buckets a day's contribution count into a 0-4 color level relative
// to the dataset's max day count, approximating GitHub's own calendar shading
// (the API does not expose an explicit per-day level).
func Level(count, max int) int {
	if count == 0 || max <= 0 {
		return 0
	}
	ratio := float64(count) / float64(max)
	switch {
	case ratio <= 0.25:
		return 1
	case ratio <= 0.5:
		return 2
	case ratio <= 0.75:
		return 3
	default:
		return 4
	}
}
