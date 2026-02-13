package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
	owner  string
	repo   string
}

func NewClient(token, owner, repo string) *Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return &Client{
		client: github.NewClient(tc),
		owner:  owner,
		repo:   repo,
	}
}

func (c *Client) SaveComment(message string) (*github.RepositoryContentResponse, error) {
	ctx := context.Background()
	fileName := fmt.Sprintf("incoming/msg-%d.md", time.Now().Unix())

	content := []byte(fmt.Sprintf("# New Feedback\n\nDate: %s\n\n%s", time.Now().Format(time.RFC1123), message))

	opts := &github.RepositoryContentFileOptions{
		Message: github.String("New anonymous feedback received"),
		Content: content,
	}

	res, _, err := c.client.Repositories.CreateFile(ctx, c.owner, c.repo, fileName, opts)
	return res, err
}
