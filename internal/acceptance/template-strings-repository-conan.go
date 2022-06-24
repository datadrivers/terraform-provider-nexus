package acceptance

const (
	TemplateStringRepositoryConanProxy = `
resource "nexus_repository_conan_proxy" "acceptance" {
` + TemplateStringProxyRepository
)
