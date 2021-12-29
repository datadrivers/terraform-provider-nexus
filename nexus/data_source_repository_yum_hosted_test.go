package nexus

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testAccDataSourceRepositoryYumHostedConfig(name string) string {
	return fmt.Sprintf(`
data "nexus_repository_yum_hosted" "acceptance" {
	name   = "%s"
}`, name)
}

func TestAccDataSourceRepositoryYumHosted(t *testing.T) {
	repo := repository.YumHostedRepository{
		Name:   fmt.Sprintf("acceptance-%s", acctest.RandString(10)),
		Online: true,
		Storage: repository.HostedStorage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: false,
		},
		Yum: repository.Yum{
			RepodataDepth: 0,
		},
	}
	nexusClient := getTestClient()
	nexusClient.Repository.Yum.Hosted.Create(repo)
	defer nexusClient.Repository.Yum.Hosted.Delete(repo.Name)
	dataSourceName := "data.nexus_repository_yum_hosted.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRepositoryYumHostedConfig(repo.Name),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "id", repo.Name),
						resource.TestCheckResourceAttr(dataSourceName, "name", repo.Name),
						resource.TestCheckResourceAttr(dataSourceName, "online", strconv.FormatBool(repo.Online)),
						resource.TestCheckResourceAttr(dataSourceName, "repodata_depth", strconv.Itoa(repo.Yum.RepodataDepth)),
						resource.TestCheckResourceAttr(dataSourceName, "storage.0.blob_store_name", repo.Storage.BlobStoreName),
						resource.TestCheckResourceAttr(dataSourceName, "storage.0.strict_content_type_validation", strconv.FormatBool(repo.Storage.StrictContentTypeValidation)),
					),
				),
			},
		},
	})
}
