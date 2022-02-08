package acceptance

const (
	TemplateStringRepositoryAptHosted = `
resource "nexus_repository_apt_hosted" "acceptance" {
	distribution = "{{ .Apt.Distribution }}"
	signing {
		keypair = "{{ .AptSigning.Keypair }}"
{{- if .AptSigning.Passphrase }}
		passphrase = "{{ .AptSigning.Passphrase }}"
{{- end }}
	}
` + TemplateStringHostedRepository

	TemplateStringRepositoryAptProxy = `
resource "nexus_repository_apt_proxy" "acceptance" {
	distribution = "{{ .Apt.Distribution }}"
	flat         = "{{ .Apt.Flat }}"

` + TemplateStringProxyRepository
)
