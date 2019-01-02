resource "github_repository_webhook" "{{ .Repository.Name}}-{{ .Hook.ID}}" {
  repository = "${github_repository.{{ .Repository.Name }}.name}"

  name = "web"

  configuration {
    {{ range $key, $value := .Hook.Config -}}
    {{ $key }} = "{{ $value }}"
    {{ end }}
  }
  active = "{{ .Hook.Active }}"
  events =  [{{ range .Hook.Events }}"{{.}}",{{end}}]
}
