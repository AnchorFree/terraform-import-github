# GitHub Terraform exporters

Tools used to migrate from manual github management to terraform based management

In order to export your organisation into Terraform, you should do:

```
export GITHUB_TOKEN="your token"
```

In order to generate terraform configuration for your members, you may need to create following template:

#### templates/admin.tpl:

```
resource "github_membership" "{{ .Login}}" {
    username = "{{ .Login}}"
    role     = "admin"
}
```

#### templates/member.tpl:

```
resource "github_membership" "{{ .Login}}" {
    username = "{{ .Login}}"
    role     = "member"
}
```

You may want to have a reference to GitHub (User)[https://godoc.org/github.com/google/go-github/github#User] for further extension of the template

After that you may need to run following command for admins export:

```
members -template templates/admin.tpl -role admin > github-admin-users.tf
terraform fmt github-admin-users.tf
```

for unprivileges members export:

```
members -template templates/member.tpl -role member > github-member-users.tf
terraform fmt github-member-users.tf
```

This will create `github-member-users.tf` and `github-admin-users.tf` files, which you can place into your terraform related folder.

In order to actually import these users into your terraform state, you may create following template:

##### templates/member_import.tpl

```
terraform import github_membership.{{ .Login }} <GitHub organization>:{{ .Login }}
```

And after that run following command:

```
members -template -fast templates/member_import.tpl > member_import.sh
```
