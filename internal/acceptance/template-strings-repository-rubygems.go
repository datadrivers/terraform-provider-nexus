package acceptance

const (
	TemplateStringRepositoryRubygemsHosted = `
resource "nexus_repository_rubygems_hosted" "acceptance" {
` + TemplateStringHostedRepository

	TemplateStringRepositoryRubygemsGroup = `
resource "nexus_repository_rubygems_group" "acceptance" {
	depends_on = [
		nexus_repository_rubygems_hosted.acceptance
	]
` + TemplateStringGroupRepository

	TemplateStringRepositoryRubygemsProxy = `
resource "nexus_repository_rubygems_proxy" "acceptance" {
` + TemplateStringProxyRepository
)
