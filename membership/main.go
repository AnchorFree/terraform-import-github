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

	tPath := flag.String("template", "templates/member.tpl", "Template to render members from")
	org := flag.String("org", "AnchorFree", "GitHub Organisation")
	publicOnly := flag.Bool("public-only", false, "list only publicly visible members, ignored if role is collaborator")
	filter := flag.String("filter", "all", "Filter members returned in the list: 2fa_disabled, all")
	role := flag.String("role", "all", "filters members by their role in organisation: all, admin, member, collaborator")

	flag.Parse()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var allMembers []*github.User
	var members []*github.User
	var resp *github.Response
	var err error

	var opt interface{}
	if *role == "collaborator" {
		opt = &github.ListOutsideCollaboratorsOptions{
			Filter:      *filter,
			ListOptions: github.ListOptions{PerPage: perPage},
		}
	} else {
		opt = &github.ListMembersOptions{
			PublicOnly:  *publicOnly,
			Filter:      *filter,
			Role:        *role,
			ListOptions: github.ListOptions{PerPage: perPage},
		}
	}

	for {
		if *role == "collaborator" {
			members, resp, err = client.Organizations.ListOutsideCollaborators(ctx, *org, opt.(*github.ListOutsideCollaboratorsOptions))
		} else {
			members, resp, err = client.Organizations.ListMembers(ctx, *org, opt.(*github.ListMembersOptions))
		}
		if err != nil {
			panic(err)
		}
		allMembers = append(allMembers, members...)
		if resp.NextPage == 0 {
			break
		}
		opt.(*github.ListMembersOptions).Page = resp.NextPage
	}

	t := template.Must(template.ParseFiles(*tPath))
	for i := range allMembers {
		err := t.Execute(os.Stdout, allMembers[i])
		if err != nil {
			panic(err)
		}
	}
}
