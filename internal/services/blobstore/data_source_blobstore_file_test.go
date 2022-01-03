package blobstore_test

import (
	"fmt"
	"testing"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceBlobstoreFile(t *testing.T) {
	bsName := "default"
	dataSourceName := "data.nexus_blobstore_file.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBlobstoreFileConfig(bsName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "default"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "default"),
					resource.TestCheckResourceAttr(dataSourceName, "path", "default"),
					resource.TestCheckResourceAttrSet(dataSourceName, "blob_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_size_in_bytes"),
					resource.TestCheckResourceAttrSet(dataSourceName, "available_space_in_bytes"),
				),
			},
		},
	})
}

func testAccDataSourceBlobstoreFileConfig(name string) string {
	return fmt.Sprintf(`
data "nexus_blobstore_file" "acceptance" {
	name = "%s"
}`, name)
}
