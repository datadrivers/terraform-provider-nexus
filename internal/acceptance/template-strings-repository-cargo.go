package acceptance

const (
	TemplateStringRepositoryCargoHosted = `
resource "nexus_repository_cargo_hosted" "acceptance" {
` + TemplateStringHostedRepository

	TemplateStringRepositoryCargoGroup = `
resource "nexus_repository_cargo_group" "acceptance" {
	depends_on = [
		nexus_repository_cargo_hosted.acceptance
	]
` + TemplateStringGroupRepository

	TemplateStringRepositoryCargoProxy = `
resource "nexus_repository_cargo_proxy" "acceptance" {
` + TemplateStringProxyRepository
)
