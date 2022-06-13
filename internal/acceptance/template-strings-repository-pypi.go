package acceptance

const (
	TemplateStringRepositoryPypiHosted = `
resource "nexus_repository_pypi_hosted" "acceptance" {
` + TemplateStringHostedRepository

	TemplateStringRepositoryPypiGroup = `
resource "nexus_repository_pypi_group" "acceptance" {
	depends_on = [
		nexus_repository_pypi_hosted.acceptance
	]
` + TemplateStringGroupRepository

	TemplateStringRepositoryPypiProxy = `
resource "nexus_repository_pypi_proxy" "acceptance" {
` + TemplateStringProxyRepository
)
