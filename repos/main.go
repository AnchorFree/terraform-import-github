package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"text/template"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
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

	// list repositories by type
	opt := &github.RepositoryListByOrgOptions{
		Type:        *repoType,
		ListOptions: github.ListOptions{PerPage: 50},
	}
	t := template.Must(template.ParseFiles(*tPath))
	for {
		repos, resp, err := client.Repositories.ListByOrg(context.Background(), *org, opt)
		if err != nil {
			fmt.Println(err)
		}
		for i, repo := range repos {
			if !*fast {
				// List doesn't provide allow_merge_commit and other parameters
				// in order to get reliable data, we need to get the data from
				// separate per repo query
				repo, _, err = client.Repositories.Get(ctx, *org, *repos[i].Name)
				if err != nil {
					panic(err)
				}
			}
			err = t.Execute(os.Stdout, repo)
			if err != nil {
				panic(err)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

}
