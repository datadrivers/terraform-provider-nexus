package acceptance

const (
	TemplateStringRepositoryNugetHosted = `
resource "nexus_repository_nuget_hosted" "acceptance" {
` + TemplateStringHostedRepository

	TemplateStringRepositoryNugetGroup = `
resource "nexus_repository_nuget_group" "acceptance" {
	depends_on = [
		nexus_repository_nuget_hosted.acceptance
	]
` + TemplateStringGroupRepository

	TemplateStringRepositoryNugetProxy = `
resource "nexus_repository_nuget_proxy" "acceptance" {
	nuget_version = "{{ .NugetProxy.NugetVersion }}"
	query_cache_item_max_age = {{ .NugetProxy.QueryCacheItemMaxAge }}
` + TemplateStringProxyRepository
)
