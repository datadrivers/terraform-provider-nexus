package repository_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccDataSourceRepositoryGitlfsHostedConfig() string {
	return `
data "nexus_repository_gitlfs_hosted" "acceptance" {
	name   = nexus_repository_gitlfs_hosted.acceptance.id
}`
}

func TestAccDataSourceRepositoryGitlfsHosted(t *testing.T) {
	repoUsingDefaults := repository.GitLfsHostedRepository{
		Name:   fmt.Sprintf("acceptance-%s", acctest.RandString(10)),
		Online: true,
		Storage: repository.HostedStorage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: false,
		},
	}
	dataSourceName := "data.nexus_repository_gitlfs_hosted.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryGitlfsHostedConfig(repoUsingDefaults) + testAccDataSourceRepositoryGitlfsHostedConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "id", repoUsingDefaults.Name),
						resource.TestCheckResourceAttr(dataSourceName, "name", repoUsingDefaults.Name),
						resource.TestCheckResourceAttr(dataSourceName, "online", strconv.FormatBool(repoUsingDefaults.Online)),
						resource.TestCheckResourceAttr(dataSourceName, "storage.0.blob_store_name", repoUsingDefaults.Storage.BlobStoreName),
						resource.TestCheckResourceAttr(dataSourceName, "storage.0.strict_content_type_validation", strconv.FormatBool(repoUsingDefaults.Storage.StrictContentTypeValidation)),
					),
				),
			},
		},
	})
}
