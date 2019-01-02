package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"text/template"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/anchorfree/github-terraform-exporters/pkg/repository"
)

type Output struct {
	Repository   *github.Repository
	Collaborator *github.User
}

const (
	perPage = 50
)

func main() {

	tPath := flag.String("template", "templates/collaborators.tpl", "Template to render collaborators from")
	org := flag.String("org", "AnchorFree", "GitHub Organisation")
	collaboratorType := flag.String("collaborator-type", "outside", "Collaborator time: outside, direct, all")
	repoType := flag.String("repo-type", "public", "Limit by repo type (public, private) if one repo not specified")
	repoName := flag.String("repo-name", "", "Collaborators for specific repo")

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

	// list collaborators by type
	opt := &github.ListCollaboratorsOptions{
		Affiliation: *collaboratorType,
		ListOptions: github.ListOptions{PerPage: 50},
	}
	t := template.Must(template.ParseFiles(*tPath))
	for repo := range repos {
		out.Repository = repo
	Collaborators:
		for {
			collabs, resp, err := client.Repositories.ListCollaborators(context.Background(), *org, *out.Repository.Name, opt)
			if err != nil {
				fmt.Println(err)
			}
			for i := range collabs {
				out.Collaborator = collabs[i]
				err = t.Execute(os.Stdout, out)
				if err != nil {
					panic(err)
				}
			}

			if resp.NextPage == 0 {
				break Collaborators
			}
			opt.Page = resp.NextPage
		}
	}
}
