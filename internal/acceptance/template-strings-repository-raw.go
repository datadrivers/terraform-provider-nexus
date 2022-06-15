package acceptance

const (
	TemplateStringRepositoryRawHosted = `
resource "nexus_repository_raw_hosted" "acceptance" {
` + TemplateStringHostedRepository

	TemplateStringRepositoryRawGroup = `
resource "nexus_repository_raw_group" "acceptance" {
	depends_on = [
		nexus_repository_raw_hosted.acceptance
	]
` + TemplateStringGroupRepository

	TemplateStringRepositoryRawProxy = `
resource "nexus_repository_raw_proxy" "acceptance" {
` + TemplateStringProxyRepository
)
