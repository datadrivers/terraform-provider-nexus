package acceptance

const (
	TemplateStringRepositoryGitlfsHosted = `
resource "nexus_repository_gitlfs_hosted" "acceptance" {
` + TemplateStringHostedRepository
)
