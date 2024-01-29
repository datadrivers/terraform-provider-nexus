package acceptance

const (
	TemplateStringRepositoryDockerHosted = `
resource "nexus_repository_docker_hosted" "acceptance" {
	docker {
		force_basic_auth = "{{ .Docker.ForceBasicAuth }}"
{{- if .Docker.HTTPPort }}
		http_port = "{{ .Docker.HTTPPort }}"
{{- end }}
{{- if .Docker.HTTPSPort }}
		https_port = "{{ .Docker.HTTPSPort }}"
{{- end }}
{{- if .Docker.Subdomain }}
		https_port = "{{ .Docker.Subdomain }}"
{{- end }}
		v1_enabled = "{{ .Docker.V1Enabled }}"
	}
` + TemplateStringHostedRepository

	TemplateStringRepositoryDockerGroup = `
resource "nexus_repository_docker_group" "acceptance" {
	docker {
		force_basic_auth = "{{ .Docker.ForceBasicAuth }}"
{{- if .Docker.HTTPPort }}
		http_port = "{{ .Docker.HTTPPort }}"
{{- end }}
{{- if .Docker.HTTPSPort }}
		https_port = "{{ .Docker.HTTPSPort }}"
{{- end }}
{{- if .Docker.Subdomain }}
		subdomain = "{{ .Docker.Subdomain }}"
{{- end }}
		v1_enabled = "{{ .Docker.V1Enabled }}"
	}
	depends_on = [
		nexus_repository_docker_hosted.acceptance
	]
` + TemplateStringGroupDeployRepository

	TemplateStringRepositoryDockerProxy = `
resource "nexus_repository_docker_proxy" "acceptance" {
	docker {
		force_basic_auth = "{{ .Docker.ForceBasicAuth }}"
{{- if .Docker.HTTPPort }}
		http_port = "{{ .Docker.HTTPPort }}"
{{- end }}
{{- if .Docker.HTTPSPort }}
		https_port = "{{ .Docker.HTTPSPort }}"
{{- end }}
{{- if .Docker.Subdomain }}
		https_port = "{{ .Docker.Subdomain }}"
{{- end }}
		v1_enabled = "{{ .Docker.V1Enabled }}"
	}

	docker_proxy {
		index_type = "{{ .DockerProxy.IndexType }}"
{{- if .DockerProxy.IndexURL }}
		index_url = "{{ .DockerProxy.IndexURL }}"
{{- end }}
	}
` + TemplateStringProxyRepository
)
