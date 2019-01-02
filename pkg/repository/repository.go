package repository

import (
	"context"

	"github.com/google/go-github/github"
)

// ListConfig configures repository listing
type ListConfig struct {
	Fast         bool
	ListOptions  github.ListOptions
	Organization string
	Type         string
	RepoName     string
}

// List repositories and send result into out channel
func List(client *github.Client, out chan *github.Repository, config ListConfig) error {
	ctx := context.Background()
	opt := &github.RepositoryListByOrgOptions{
		Type:        config.Type,
		ListOptions: config.ListOptions,
	}

	if config.RepoName != "" {
		repo, _, err := client.Repositories.Get(ctx, config.Organization, config.RepoName)
		if err != nil {
			return err
		}
		out <- repo
		return nil
	}

	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, config.Organization, opt)
		if err != nil {
			return err
		}
		for i := range repos {
			if config.Fast {
				out <- repos[i]
			} else {
				repo, _, err := client.Repositories.Get(ctx, config.Organization, *repos[i].Name)
				if err != nil {
					return err
				}
				out <- repo
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return nil
}
