package acceptance

const (
	TemplateStringRepositoryHelmHosted = `
resource "nexus_repository_helm_hosted" "acceptance" {
` + TemplateStringHostedRepository

	TemplateStringRepositoryHelmProxy = `
resource "nexus_repository_helm_proxy" "acceptance" {
` + TemplateStringProxyRepository
)
