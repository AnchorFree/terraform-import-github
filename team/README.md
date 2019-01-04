# GitHub Terraform exporters

Tools used to migrate from manual github management to terraform based management

In order to export your organisation into Terraform, you should do:

```
export GITHUB_TOKEN="your token"
```

In order to generate terraform configuration for your teams, you may need to create following template:

#### templates/team.tpl:

```
resource "github_team" "{{ .Slug }}" {
    name = "{{ .Name }}"
    description = "{{ .Description }}"
    privacy = "{{ .Privacy }}"
    {{ if .Parent -}}
    parent_team_id = {{ .Parent.ID }}
{{ end -}}
}

```

You may want to have a reference to GitHub (Team)[https://godoc.org/github.com/google/go-github/github#Team] for further extension of the template

After that you may need to run following command for teams export:

```
team -template templates/team.tpl > github-teams.tf
terraform fmt github-teams.tf
```

This will create `github-team.tf` file, which you can place into your terraform related folder.

In order to actually import these teams into your terraform state, you may create following template:

##### templates/team_import.tpl

```
terraform import github_team.{{ .Slug }} {{ .ID }}
```

And after that run following command:

```
members -template templates/team_import.tpl > team_import.sh
```
