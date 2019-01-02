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
