package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"strings"
)

// Client wraps up github.Client
type Client struct {
	github.Client
}

func NewGithubClient() Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "... your access token ..."},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return Client{*client}
}

// getReadmeLines returns slice of string each representing a line in avelino/awesome-go project
func (c Client) getReadmeLines(ctx context.Context) []string {
	fc, _, _, err := c.Repositories.GetContents(ctx, "avelino", "awesome-go", "README.md", nil)
	if err != nil {
		fmt.Println("error:", err)
	}
	data, err := base64.StdEncoding.DecodeString(*fc.Content)
	if err != nil {
		fmt.Println("error:", err)
	}

	return strings.Split(strings.Replace(string(data), "\r\n", "\n", -1), "\n")
}

// fetchRepository returns a repository
func (c Client) fetchRepository(ctx context.Context, owner, repo string) (*github.Repository, error) {
	repository, _, err := c.Repositories.Get(ctx, owner, repo)
	return repository, err
}
