package acceptance

const (
	TemplateStringRepositoryBowerHosted = `
resource "nexus_repository_bower_hosted" "acceptance" {
` + TemplateStringHostedRepository

	TemplateStringRepositoryBowerGroup = `
resource "nexus_repository_bower_group" "acceptance" {
	depends_on = [
		nexus_repository_bower_hosted.acceptance
	]
` + TemplateStringGroupRepository

	TemplateStringRepositoryBowerProxy = `
resource "nexus_repository_bower_proxy" "acceptance" {
	rewrite_package_urls = {{ .Bower.RewritePackageUrls }}
` + TemplateStringProxyRepository
)
