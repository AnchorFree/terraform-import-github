package main

import (
	"context"
	"flag"
	"os"

	"text/template"

	"github.com/anchorfree/github-terraform-exporters/pkg/repository"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	perPage = 50
)

func main() {

	tPath := flag.String("template", "templates/repo.tpl", "Template to render repos from")
	org := flag.String("org", "AnchorFree", "GitHub Organisation")
	repoType := flag.String("type", "public", "Limit by repo type (public, private)")
	fast := flag.Bool("fast", false, "Don't run per repo additional query, some parameters are not passed otherwise")

	flag.Parse()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repoConfig := repository.ListConfig{
		ListOptions:  github.ListOptions{PerPage: perPage},
		Type:         *repoType,
		Organization: *org,
		Fast:         *fast,
	}

	repos := make(chan *github.Repository, perPage)
	go func() {
		err := repository.List(client, repos, repoConfig)
		if err != nil {
			panic(err)
		}
		close(repos)
	}()

	t := template.Must(template.ParseFiles(*tPath))
	for repo := range repos {
		err := t.Execute(os.Stdout, repo)
		if err != nil {
			panic(err)
		}
	}
}
