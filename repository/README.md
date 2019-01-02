# GitHub Terraform exporters

Tools used to migrate from manual github management to terraform based management

In order to export your organisation into Terraform, you should do:

```
export GITHUB_TOKEN="your token"
```

In order to generate terraform configuration for your public repos, you may need to create following template:

#### templates/repo.tpl:

```
resource "github_repository" "{{ .Name }}" {
    name        = "{{ .Name }}"
    private     = "{{ .Private }}"
    {{ if .Description }}
    description = "{{ .Description }}"
    {{ end }}
    has_wiki    = "{{ .HasWiki }}"
    has_downloads = "{{ .HasDownloads }}"
    has_issues  = "{{ .HasIssues }}"
    has_projects  = "{{ .HasProjects }}"
    {{ if .Homepage }}
    homepage_url = "{{ .Homepage }}"
    {{ end }}
    allow_merge_commit  = "{{ .AllowMergeCommit }}"
    allow_squash_merge  = "{{ .AllowSquashMerge }}"
    allow_rebase_merge  = "{{ .AllowRebaseMerge }}"
    {{ if .Topics }}
    topics = [{{ range .Topics }} "{{.}}", {{end}}]
    {{ end }}
    archived = "{{ .Archived }}"
}
```

You may want to have a reference to GitHub (Repository)[https://godoc.org/github.com/google/go-github/github#Repository] for further extension of the template

After that you may need to run following command:

```
repos -template templates/repo.tpl -type public > github-public-repos.tf
terraform fmt github-public-repos.tf
```

This will create `github-public-repos.tf` file, which you can place into your terraform related folder.

In order to actually import these repos into your terraform state, you may create following template:

##### templates/repo_import.tpl

```
terraform import github_repository.{{ .Name }} {{ .Name }}
```

And after that run following command:

```
repos -template -fast templates/repo_import.tpl -type public > public_import.sh
```

You may noticed `fast` argument, it will use repo `List` API instead of `Get`, which is several times faster, but doesn't provide complete information about repo, such as:

- `allow_squash_merge`
- `allow_merge_commit`
- `allow_rebase_merge`

If you don't use those fields in your templates - you may want to skip it.
