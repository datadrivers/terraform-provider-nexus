package acceptance

const (
	TemplateStringRepositoryRawHosted = `
resource "nexus_repository_raw_hosted" "{{ .Name }}" {
` + TemplateStringHostedRepository

	TemplateStringRepositoryRawGroup = `
resource "nexus_repository_raw_group" "acceptance" {
	depends_on = [
	{{- range $val := .Group.MemberNames }}
		nexus_repository_raw_hosted.{{ $val }},
	{{ end -}}
	]
` + TemplateStringGroupRepository

	TemplateStringRepositoryRawProxy = `
resource "nexus_repository_raw_proxy" "acceptance" {
` + TemplateStringProxyRepository
)
