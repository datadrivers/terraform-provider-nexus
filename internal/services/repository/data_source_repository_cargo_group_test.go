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

func testAccDataSourceRepositoryCargoGroupConfig() string {
	return `
data "nexus_repository_cargo_group" "acceptance" {
	name   = nexus_repository_cargo_group.acceptance.id
}`
}

func TestAccDataSourceRepositoryCargoGroup(t *testing.T) {
	repoHosted := testAccResourceRepositoryCargoHosted()
	repoGroup := repository.CargoGroupRepository{
		Name:   fmt.Sprintf("acceptance-%s", acctest.RandString(10)),
		Online: true,
		Storage: repository.Storage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: false,
		},
		Group: repository.Group{
			MemberNames: []string{repoHosted.Name},
		},
	}
	dataSourceName := "data.nexus_repository_cargo_group.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryCargoHostedConfig(repoHosted) + testAccResourceRepositoryCargoGroupConfig(repoGroup) + testAccDataSourceRepositoryCargoGroupConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "id", repoGroup.Name),
							resource.TestCheckResourceAttr(dataSourceName, "name", repoGroup.Name),
							resource.TestCheckResourceAttr(dataSourceName, "online", strconv.FormatBool(repoGroup.Online)),
						),
						resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "storage.#", "1"),
							resource.TestCheckResourceAttr(dataSourceName, "storage.0.blob_store_name", repoGroup.Storage.BlobStoreName),
							resource.TestCheckResourceAttr(dataSourceName, "storage.0.strict_content_type_validation", strconv.FormatBool(repoGroup.Storage.StrictContentTypeValidation)),
							resource.TestCheckResourceAttr(dataSourceName, "group.#", "1"),
							resource.TestCheckResourceAttr(dataSourceName, "group.0.member_names.#", "1"),
							resource.TestCheckResourceAttr(dataSourceName, "group.0.member_names.0", repoGroup.Group.MemberNames[0]),
						),
					),
				),
			},
		},
	})
}
