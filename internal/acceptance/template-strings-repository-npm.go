package acceptance

const (
	TemplateStringRepositoryNpmHosted = `
resource "nexus_repository_npm_hosted" "acceptance" {
` + TemplateStringHostedRepository

	TemplateStringRepositoryNpmGroup = `
resource "nexus_repository_npm_group" "acceptance" {
	depends_on = [
		nexus_repository_npm_hosted.acceptance
	]
` + TemplateStringGroupRepository

	TemplateStringRepositoryNpmProxy = `
resource "nexus_repository_npm_proxy" "acceptance" {
	remove_quarantined = {{ .Npm.RemoveQuarantined }}
	remove_non_cataloged = {{ .Npm.RemoveNonCataloged }}
` + TemplateStringProxyRepository
)
