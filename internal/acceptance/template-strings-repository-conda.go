package acceptance

const (
	TemplateStringRepositoryCondaProxy = `
resource "nexus_repository_conda_proxy" "acceptance" {
` + TemplateStringProxyRepository
)
