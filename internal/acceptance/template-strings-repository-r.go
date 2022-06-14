package acceptance

const (
	TemplateStringRepositoryRHosted = `
resource "nexus_repository_r_hosted" "acceptance" {
` + TemplateStringHostedRepository

	TemplateStringRepositoryRGroup = `
resource "nexus_repository_r_group" "acceptance" {
	depends_on = [
		nexus_repository_r_hosted.acceptance
	]
` + TemplateStringGroupRepository

	TemplateStringRepositoryRProxy = `
resource "nexus_repository_r_proxy" "acceptance" {
` + TemplateStringProxyRepository
)
