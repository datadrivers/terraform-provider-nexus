package acceptance

const (
	TemplateStringRepositoryMavenHosted = `
resource "nexus_repository_maven_hosted" "acceptance" {
	maven {
		version_policy = "{{ .Maven.VersionPolicy }}"
		layout_policy = "{{ .Maven.LayoutPolicy }}"
{{- if .Maven.ContentDisposition }}
		content_disposition = "{{ .Maven.ContentDisposition }}"
{{- end }}
	}
` + TemplateStringHostedRepository

	TemplateStringRepositoryMavenGroup = `
resource "nexus_repository_maven_group" "acceptance" {
	depends_on = [
		nexus_repository_maven_hosted.acceptance
	]
` + TemplateStringGroupRepository

	TemplateStringRepositoryMavenProxy = `
resource "nexus_repository_maven_proxy" "acceptance" {
	maven {
		version_policy = "{{ .Maven.VersionPolicy }}"
		layout_policy = "{{ .Maven.LayoutPolicy }}"
{{- if .Maven.ContentDisposition }}
		content_disposition = "{{ .Maven.ContentDisposition }}"
{{- end }}
	}
` + TemplateStringProxyRepository
)
