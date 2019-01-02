# GitHub Terraform exporters

Tools used to migrate from manual github management to terraform based management

In order to export your organisation into Terraform, you should do:

```
export GITHUB_TOKEN="your token"
```

In order to generate terraform configuration for collaborators of your repos, you may need to create following template:

#### templates/collaborators.tpl:

```
{{ $permission := "pull" -}}
{{ range $key, $value := .Collaborator.Permissions -}}
{{ if $value -}}
  {{ if eq $key "admin" -}}
    {{ $permission = $key -}}
{{ else if and (eq $key "push") (ne $permission "admin") -}}
    {{ $permission = $key -}}
{{ end -}}
{{ end -}}
{{ end -}}
resource "github_repository_collaborator" "{{ .Repository.Name}}-{{ .Collaborator.Login }}" {
    repository = "{{ .Repository.Name }}"
    username   = "{{ .Collaborator.Login }}"
    permission = "{{ $permission }}"
}
```

You may want to have a reference to GitHub (User)[https://godoc.org/github.com/google/go-github/github#User] and (Repository)[https://godoc.org/github.com/google/go-github/github#Repository]

After that you may need to run following command:

```
repo_collaborators -template templates/collaborators.tpl -repo-type public -collaborator-type outside > github-public-repo-collaborators.tf
terraform fmt github-public-repo-collaborators.tf
```

This will create `github-public-repo-collaborators.tf` file, which you can place into your terraform related folder.

In order to actually import these repo collaborators into your terraform state, you may create following template:

##### templates/collaborators_import.tpl

```
terraform import github_repository_collaborator.{{ .Repository.Name}}-{{ .Collaborator.Login }} {{ .Repository.Name}}-{{ .Collaborator.Login }}
```

And after that run following command:

```
repo_collaborators -template  templates/collaborator_import.tpl -repo-type public -collaborator-type outside  > collaborator_import.sh
```

The next step is to make it more manageable, you may want to create some terraform module, which would do loop over your push, pull, admin users per repo.

Creation of the
