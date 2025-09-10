package yt

import (
	"context"
	"encoding/json"
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
	searchURL, err := url.JoinPath(yc.baseURL, "api/issues")
	if err != nil {
		return service.Issue{}, errors.WithStack(err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return service.Issue{}, errors.WithStack(err)
	}

	q := req.URL.Query()
	q.Add("query", issueID)
	q.Add("fields", "idReadable,summary,description,tags(name)")
	req.URL.RawQuery = q.Encode()

	yc.addAuthHeaders(req)

	resp, err := yc.httpClient.Do(req)

	if err != nil {
		return service.Issue{}, errors.WithStack(err)
	}
	defer resp.Body.Close()

	var issues []youtrackIssueResponse
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return service.Issue{}, errors.WithStack(err)
	}

	if len(issues) == 0 {
		return service.Issue{}, service.ErrIssueNotFound
	}

	return yc.convertToIssue(issues[0]), nil
}

func (yc *YoutrackClient) convertToIssue(resp youtrackIssueResponse) service.Issue {
	issue := service.Issue{
		ID:          resp.ID,
		Title:       resp.Summary,
		Description: resp.Description,
		Tags:        make([]string, 0, len(resp.Tags)),
	}

	for _, tag := range resp.Tags {
		issue.Tags = append(issue.Tags, tag.Name)
	}

	return issue
}

type youtrackIssueResponse struct {
	ID          string `json:"idReadable"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Tags        []struct {
		Name string `json:"name"`
	} `json:"tags"`
}
