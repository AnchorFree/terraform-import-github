package main

import (
	"context"
	"flag"
	"os"

	"text/template"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	perPage = 50
)

func main() {

	tPath := flag.String("template", "templates/team.tpl", "Template to render members from")
	org := flag.String("org", "AnchorFree", "GitHub Organisation")

	flag.Parse()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var allTeams []*github.Team
	opt := &github.ListOptions{PerPage: perPage}

	for {
		teams, resp, err := client.Teams.ListTeams(ctx, *org, opt)
		if err != nil {
			panic(err)
		}
		allTeams = append(allTeams, teams...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	t := template.Must(template.ParseFiles(*tPath))
	for i := range allTeams {
		err := t.Execute(os.Stdout, allTeams[i])
		if err != nil {
			panic(err)
		}
	}
}
