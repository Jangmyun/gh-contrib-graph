package github

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/cli/go-gh/v2/pkg/api"
)

// roundTripFunc lets a plain function satisfy http.RoundTripper, so tests
// can fake API responses without touching the network or gh's local config.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// jsonResponse builds a canned response for req. The Request field must be
// set for go-gh's error handling path (HandleHTTPError reads resp.Request.URL).
func jsonResponse(req *http.Request, status int, body string) *http.Response {
	return &http.Response{
		Request:    req,
		StatusCode: status,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

// fakeClientOptions builds ClientOptions that skip gh's local config/auth
// resolution entirely (see optionsNeedResolution in go-gh), so these tests
// run the same with or without a `gh auth login` session available.
func fakeClientOptions(transport http.RoundTripper) api.ClientOptions {
	return api.ClientOptions{
		Host:      "github.com",
		AuthToken: "test-token",
		Transport: transport,
	}
}

func TestCurrentLogin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		rest, err := api.NewRESTClient(fakeClientOptions(roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if !strings.HasSuffix(r.URL.Path, "/user") {
				t.Errorf("unexpected request path: %s", r.URL.Path)
			}
			return jsonResponse(r, 200, `{"login":"octocat"}`), nil
		})))
		if err != nil {
			t.Fatalf("building REST client: %v", err)
		}

		login, err := CurrentLogin(rest)
		if err != nil {
			t.Fatalf("CurrentLogin returned error: %v", err)
		}
		if login != "octocat" {
			t.Errorf("CurrentLogin = %q, want %q", login, "octocat")
		}
	})

	t.Run("http error", func(t *testing.T) {
		rest, err := api.NewRESTClient(fakeClientOptions(roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return jsonResponse(r, 401, `{"message":"Bad credentials"}`), nil
		})))
		if err != nil {
			t.Fatalf("building REST client: %v", err)
		}

		if _, err := CurrentLogin(rest); err == nil {
			t.Error("expected an error for a 401 response, got nil")
		}
	})
}

func TestFetchContributions(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gql, err := api.NewGraphQLClient(fakeClientOptions(roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return jsonResponse(r, 200, `{
				"data": {
					"user": {
						"contributionsCollection": {
							"contributionCalendar": {
								"totalContributions": 42,
								"weeks": [
									{"contributionDays": [{"date":"2026-01-04","contributionCount":3,"weekday":0}]}
								]
							}
						}
					}
				}
			}`), nil
		})))
		if err != nil {
			t.Fatalf("building GraphQL client: %v", err)
		}

		cal, err := FetchContributions(gql, "octocat")
		if err != nil {
			t.Fatalf("FetchContributions returned error: %v", err)
		}
		if cal.TotalContributions != 42 {
			t.Errorf("TotalContributions = %d, want 42", cal.TotalContributions)
		}
		if len(cal.Weeks) != 1 || len(cal.Weeks[0].ContributionDays) != 1 {
			t.Fatalf("unexpected weeks shape: %+v", cal.Weeks)
		}
		if got := cal.Weeks[0].ContributionDays[0].ContributionCount; got != 3 {
			t.Errorf("ContributionCount = %d, want 3", got)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		gql, err := api.NewGraphQLClient(fakeClientOptions(roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return jsonResponse(r, 200, `{"data": {"user": null}}`), nil
		})))
		if err != nil {
			t.Fatalf("building GraphQL client: %v", err)
		}

		if _, err := FetchContributions(gql, "does-not-exist"); err == nil {
			t.Error("expected an error for a nil user, got nil")
		}
	})

	t.Run("graphql error", func(t *testing.T) {
		gql, err := api.NewGraphQLClient(fakeClientOptions(roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return jsonResponse(r, 200, `{"errors":[{"message":"Could not resolve to a User"}]}`), nil
		})))
		if err != nil {
			t.Fatalf("building GraphQL client: %v", err)
		}

		if _, err := FetchContributions(gql, "does-not-exist"); err == nil {
			t.Error("expected an error when the API returns a GraphQL error, got nil")
		}
	})
}
