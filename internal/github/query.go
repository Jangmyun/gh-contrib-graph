package github

// ContributionDay is a single day's contribution count within a calendar week.
type ContributionDay struct {
	Date              string `json:"date"`
	ContributionCount int    `json:"contributionCount"`
	Weekday           int    `json:"weekday"`
}

// ContributionWeek groups the days GitHub renders as one column in the calendar.
type ContributionWeek struct {
	ContributionDays []ContributionDay `json:"contributionDays"`
}

// ContributionCalendar is the rolling-year contribution calendar for a user.
type ContributionCalendar struct {
	TotalContributions int                `json:"totalContributions"`
	Weeks              []ContributionWeek `json:"weeks"`
}

type contributionsCollection struct {
	ContributionCalendar ContributionCalendar `json:"contributionCalendar"`
}

type userResponse struct {
	ContributionsCollection contributionsCollection `json:"contributionsCollection"`
}

type contributionQueryResponse struct {
	User *userResponse `json:"user"`
}

const contributionQuery = `
query($login: String!, $from: DateTime!, $to: DateTime!) {
  user(login: $login) {
    contributionsCollection(from: $from, to: $to) {
      contributionCalendar {
        totalContributions
        weeks {
          contributionDays {
            date
            contributionCount
            weekday
          }
        }
      }
    }
  }
}`
