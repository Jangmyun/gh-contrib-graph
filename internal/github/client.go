package github

import (
	"fmt"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
)

// NewGraphQLClient builds a GraphQL client that reuses the local `gh auth` session.
func NewGraphQLClient() (*api.GraphQLClient, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("building GitHub GraphQL client (is `gh auth login` set up?): %w", err)
	}
	return client, nil
}

// NewRESTClient builds a REST client that reuses the local `gh auth` session.
func NewRESTClient() (*api.RESTClient, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("building GitHub REST client (is `gh auth login` set up?): %w", err)
	}
	return client, nil
}

// CurrentLogin resolves the login of the currently authenticated gh user.
func CurrentLogin(rest *api.RESTClient) (string, error) {
	var resp struct {
		Login string `json:"login"`
	}
	if err := rest.Get("user", &resp); err != nil {
		return "", fmt.Errorf("resolving authenticated user: %w", err)
	}
	return resp.Login, nil
}

// FetchContributions fetches the rolling-year contribution calendar for login.
func FetchContributions(client *api.GraphQLClient, login string) (*ContributionCalendar, error) {
	now := time.Now()
	vars := map[string]interface{}{
		"login": login,
		"from":  now.AddDate(-1, 0, 0).Format(time.RFC3339),
		"to":    now.Format(time.RFC3339),
	}

	var resp contributionQueryResponse
	if err := client.Do(contributionQuery, vars, &resp); err != nil {
		return nil, fmt.Errorf("fetching contributions for %s: %w", login, err)
	}
	if resp.User == nil {
		return nil, fmt.Errorf("user %q not found", login)
	}
	return &resp.User.ContributionsCollection.ContributionCalendar, nil
}
