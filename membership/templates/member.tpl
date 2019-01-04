resource "github_membership" "{{ .Login}}" {
    username = "{{ .Login}}"
    role     = "member"
}
