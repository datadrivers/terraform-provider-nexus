package acceptance

const (
	TemplateStringRepositoryTerraformHosted = `
resource "nexus_repository_terraform_hosted" "{{ .Name }}" {
` + TemplateStringHostedRepository

	TemplateStringRepositoryTerraformGroup = `
resource "nexus_repository_terraform_group" "acceptance" {
	depends_on = [
	{{- range $val := .Group.MemberNames }}
		nexus_repository_terraform_hosted.{{ $val }},
	{{ end -}}
	]
` + TemplateStringGroupRepository

	TemplateStringRepositoryTerraformProxy = `
resource "nexus_repository_terraform_proxy" "acceptance" {
` + TemplateStringProxyRepository
)
