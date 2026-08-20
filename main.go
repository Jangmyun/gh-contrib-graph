package main

import (
	"fmt"
	"os"

	"github.com/Jangmyun/gh-contrib-graph/internal/contrib"
	ghapi "github.com/Jangmyun/gh-contrib-graph/internal/github"
	"github.com/Jangmyun/gh-contrib-graph/internal/render"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "gh-contrib-graph:", err)
		os.Exit(1)
	}
}

func run() error {
	login := ""
	if len(os.Args) > 1 {
		login = os.Args[1]
	}

	rest, err := ghapi.NewRESTClient()
	if err != nil {
		return err
	}
	if login == "" {
		login, err = ghapi.CurrentLogin(rest)
		if err != nil {
			return err
		}
	}

	graphql, err := ghapi.NewGraphQLClient()
	if err != nil {
		return err
	}

	cal, err := ghapi.FetchContributions(graphql, login)
	if err != nil {
		return err
	}

	var days []ghapi.ContributionDay
	for _, w := range cal.Weeks {
		days = append(days, w.ContributionDays...)
	}
	longest := contrib.LongestStreak(days)
	current := contrib.CurrentStreak(days)

	fmt.Println(render.Render(login, cal, longest, current))
	return nil
}
