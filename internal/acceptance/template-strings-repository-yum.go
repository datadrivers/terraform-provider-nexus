package acceptance

const (
	TemplateStringRepositoryYumHosted = `
resource "nexus_repository_yum_hosted" "acceptance" {
	{{- if .Yum.DeployPolicy }}
	deploy_policy  = "{{ .Yum.DeployPolicy }}"
	{{- end }}
	repodata_depth = {{ .Yum.RepodataDepth }}
` + TemplateStringHostedRepository

	TemplateStringRepositoryYumGroup = `
resource "nexus_repository_yum_group" "acceptance" {
{{ if .YumSigning }}
	yum_signing {
		keypair = "{{ .YumSigning.Keypair }}"
{{ if .YumSigning.Passphrase }}
		passphrase = "{{ .YumSigning.Passphrase }}"
{{ end -}}
	}
{{ end -}}
	depends_on = [
		nexus_repository_yum_hosted.acceptance
	]
` + TemplateStringGroupRepository

	TemplateStringRepositoryYumProxy = `
resource "nexus_repository_yum_proxy" "acceptance" {
{{ if .YumSigning }}
	yum_signing {
		keypair = "{{ .YumSigning.Keypair }}"
{{ if .YumSigning.Passphrase }}
		passphrase = "{{ .YumSigning.Passphrase }}"
{{ end -}}
	}
{{ end -}}
` + TemplateStringProxyRepository
)
