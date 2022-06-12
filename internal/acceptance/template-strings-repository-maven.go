package acceptance

const (
	TemplateStringRepositoryMavenHosted = `
resource "nexus_repository_maven_hosted" "acceptance" {
	maven {
{{- if .Maven.VersionPolicy }}
		version_policy = "{{ .Maven.VersionPolicy }}"
{{- end }}
{{- if .Maven.LayoutPolicy }}
		layout_policy = "{{ .Maven.LayoutPolicy }}"
{{- end }}
{{- if .Maven.ContentDisposition }}
		content_disposition = "{{ .Maven.ContentDisposition }}"
{{- end }}
	}
` + TemplateStringHostedRepository
)
