package yt

import (
	"context"
	"encoding/json"
	"fmt"
	"gitrack/pkg/app/service"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type YoutrackClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewService(baseURL, token string) *YoutrackClient {
	return &YoutrackClient{
		baseURL:    baseURL,
		token:      token,
		httpClient: &http.Client{},
	}
}

func (yc *YoutrackClient) addAuthHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+yc.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
}

func (yc *YoutrackClient) GetIssue(ctx context.Context, issueID string) (service.Issue, error) {
	issueURL, err := url.JoinPath(yc.baseURL, "api/issues", issueID)
	if err != nil {
		return service.Issue{}, errors.WithStack(err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", issueURL, nil)
	if err != nil {
		return service.Issue{}, errors.WithStack(err)
	}

	yc.addAuthHeaders(req)

	resp, err := yc.httpClient.Do(req)
	if err != nil {
		return service.Issue{}, errors.WithStack(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return service.Issue{}, errors.WithStack(service.ErrIssueNotFound)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return service.Issue{}, errors.WithStack(fmt.Errorf("authentication failed: invalid token"))
	}
	if resp.StatusCode != http.StatusOK {
		return service.Issue{}, errors.WithStack(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
	}

	var issueResponse youtrackIssueResponse
	if err := json.NewDecoder(resp.Body).Decode(&issueResponse); err != nil {
		return service.Issue{}, errors.WithStack(err)
	}

	return yc.convertToIssue(issueResponse), nil
}

func (yc *YoutrackClient) convertToIssue(resp youtrackIssueResponse) service.Issue {
	issue := service.Issue{
		ID:          resp.ID,
		Title:       resp.Summary,
		Description: resp.Description,
		Tags:        make([]string, 0, len(resp.Tags)),
	}

	switch resp.State.Name {
	case "Code Review":
		issue.State = service.IssueStateCodeReview
	default:
		issue.State = service.IssueStateOther
	}

	for _, tag := range resp.Tags {
		issue.Tags = append(issue.Tags, tag.Name)
	}

	return issue
}

type youtrackIssueResponse struct {
	ID          string `json:"id"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	State       struct {
		Type string `json:"$type"`
		Name string `json:"name"`
	} `json:"state"`
	Tags []struct {
		Name string `json:"name"`
	} `json:"tags"`
}
