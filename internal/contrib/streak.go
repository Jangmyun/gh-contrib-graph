package contrib

import "github.com/Jangmyun/gh-contrib-graph/internal/github"

// LongestStreak returns the longest run of consecutive days with at least
// one contribution, over the chronologically ordered days slice.
func LongestStreak(days []github.ContributionDay) int {
	longest, current := 0, 0
	for _, d := range days {
		if d.ContributionCount > 0 {
			current++
			if current > longest {
				longest = current
			}
		} else {
			current = 0
		}
	}
	return longest
}

// CurrentStreak returns the length of the still-active streak ending at the
// most recent day. A trailing zero-count day (i.e. "today", which may not
// have happened yet) does not by itself break the streak; any other
// zero-count day does.
func CurrentStreak(days []github.ContributionDay) int {
	n := len(days)
	if n == 0 {
		return 0
	}

	i := n - 1
	if days[i].ContributionCount == 0 {
		i--
	}

	streak := 0
	for ; i >= 0; i-- {
		if days[i].ContributionCount == 0 {
			break
		}
		streak++
	}
	return streak
}
