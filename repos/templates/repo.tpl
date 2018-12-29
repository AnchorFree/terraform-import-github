resource "github_repository" "{{ .Name }}" {
    name        = "{{ .Name }}"
    private     = "{{ .Private }}"
    {{ if .Description -}}
    description = "{{ .Description }}"
    {{ end -}}
    has_wiki    = "{{ .HasWiki }}"
    has_downloads = "{{ .HasDownloads }}"
    has_issues  = "{{ .HasIssues }}"
    has_projects  = "{{ .HasProjects }}"
    {{ if .Homepage -}}
    homepage_url = "{{ .Homepage }}"
    {{ end -}}
    allow_merge_commit  = "{{ .AllowMergeCommit }}"
    allow_squash_merge  = "{{ .AllowSquashMerge }}"
    allow_rebase_merge  = "{{ .AllowRebaseMerge }}"
    {{ if .Topics -}}
    topics = [{{ range .Topics }} "{{.}}", {{end}}]
    {{ end -}}
    archived = "{{ .Archived }}"
}
