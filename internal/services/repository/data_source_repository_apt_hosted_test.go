package repository_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccDataSourceRepositoryAptHostedConfig() string {
	return `
data "nexus_repository_apt_hosted" "acceptance" {
	name   = nexus_repository_apt_hosted.acceptance.id
}`
}

func TestAccDataSourceRepositoryAptHosted(t *testing.T) {

	repo := repository.AptHostedRepository{
		Name:   fmt.Sprintf("acceptance-%s", acctest.RandString(10)),
		Online: true,
		Storage: repository.HostedStorage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: false,
		},
		Apt: repository.AptHosted{
			Distribution: "bullseye",
		},
		AptSigning: repository.AptSigning{
			Keypair:    "test-data-keypair",
			Passphrase: tools.GetStringPointer("test-data-passphrase"),
		},
	}
	dataSourceName := "data.nexus_repository_apt_hosted.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryAptHostedConfig(repo) + testAccDataSourceRepositoryAptHostedConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "id", repo.Name),
						resource.TestCheckResourceAttr(dataSourceName, "name", repo.Name),
						resource.TestCheckResourceAttr(dataSourceName, "online", strconv.FormatBool(repo.Online)),
						resource.TestCheckResourceAttr(dataSourceName, "distribution", repo.Apt.Distribution),
						resource.TestCheckResourceAttr(dataSourceName, "storage.0.blob_store_name", repo.Storage.BlobStoreName),
						resource.TestCheckResourceAttr(dataSourceName, "storage.0.strict_content_type_validation", strconv.FormatBool(repo.Storage.StrictContentTypeValidation)),
						resource.TestCheckResourceAttr(dataSourceName, "signing.0.keypair", repo.AptSigning.Keypair),
					),
				),
			},
		},
	})
}
