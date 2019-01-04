resource "github_team" "{{ .Slug }}" {
    name = "{{ .Name }}"
    description = "{{ .Description }}"
    privacy = "{{ .Privacy }}"
    {{ if .Parent -}}
    parent_team_id = {{ .Parent.ID }}
{{ end -}}
}
