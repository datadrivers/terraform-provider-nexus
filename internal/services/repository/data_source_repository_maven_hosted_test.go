package repository_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccDataSourceRepositoryMavenHostedConfig() string {
	return `
data "nexus_repository_maven_hosted" "acceptance" {
	name   = nexus_repository_maven_hosted.acceptance.id
}`
}

func TestAccDataSourceRepositoryMavenHosted(t *testing.T) {
	repo := repository.MavenHostedRepository{
		Name:   fmt.Sprintf("acceptance-%s", acctest.RandString(10)),
		Online: true,
		Storage: repository.HostedStorage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: false,
		},
		Maven: repository.Maven{
			VersionPolicy: repository.MavenVersionPolicyRelease,
			LayoutPolicy:  repository.MavenLayoutPolicyStrict,
		},
	}
	dataSourceName := "data.nexus_repository_maven_hosted.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryMavenHostedConfig(repo) + testAccDataSourceRepositoryMavenHostedConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "id", repo.Name),
						resource.TestCheckResourceAttr(dataSourceName, "name", repo.Name),
						resource.TestCheckResourceAttr(dataSourceName, "online", strconv.FormatBool(repo.Online)),
						resource.TestCheckResourceAttr(dataSourceName, "storage.0.blob_store_name", repo.Storage.BlobStoreName),
						resource.TestCheckResourceAttr(dataSourceName, "storage.0.strict_content_type_validation", strconv.FormatBool(repo.Storage.StrictContentTypeValidation)),
					),
				),
			},
		},
	})
}
