package main

import (
	"context"
	"flag"
	"os"
	"regexp"

	"text/template"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/anchorfree/github-terraform-exporters/pkg/repository"
)

type Output struct {
	Repository *github.Repository
	Hook       *github.Hook
}

const (
	perPage = 50
)

func main() {

	tPath := flag.String("template", "templates/webhook.tpl", "Template to render webhooks from")
	org := flag.String("org", "AnchorFree", "GitHub Organisation")
	repoType := flag.String("repo-type", "public", "Limit by repo type (public, private) if one repo not specified")
	repoName := flag.String("repo-name", "", "Collaborators for specific repo")
	filterOut := flag.String("filter-out", "", "Skip WebHooks which URL is like regexp")

	flag.Parse()

	out := new(Output)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repos := make(chan *github.Repository, perPage)

	repoConfig := repository.ListConfig{
		ListOptions:  github.ListOptions{PerPage: perPage},
		Type:         *repoType,
		Organization: *org,
		Fast:         true,
		RepoName:     *repoName,
	}

	go func() {
		err := repository.List(client, repos, repoConfig)
		if err != nil {
			panic(err)
		}
		close(repos)
	}()

	skip, _ := regexp.Compile(*filterOut)
	// list WebHooks
	opt := &github.ListOptions{PerPage: 50}
	t := template.Must(template.ParseFiles(*tPath))
	for repo := range repos {
		out.Repository = repo
	Hooks:
		for {
			hooks, resp, err := client.Repositories.ListHooks(context.Background(), *org, *out.Repository.Name, opt)
			if err != nil {
				panic(err)
			}
		hook:
			for i := range hooks {
				out.Hook = hooks[i]
				if out.Hook.Config["url"] != nil && skip.MatchString(out.Hook.Config["url"].(string)) {
					continue hook
				}
				err = t.Execute(os.Stdout, out)
				if err != nil {
					panic(err)
				}
			}

			if resp.NextPage == 0 {
				break Hooks
			}
			opt.Page = resp.NextPage
		}
	}
}
