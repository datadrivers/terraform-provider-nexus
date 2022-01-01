package blobstore_test

import (
	"fmt"
	"testing"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceBlobstoreFile(t *testing.T) {
	bsName := "default"
	resourceName := "data.nexus_blobstore_file.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBlobstoreFileConfig(bsName),
				Check: resource.ComposeTestCheckFunc(
					// Base and common resource props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", "default"),
						resource.TestCheckResourceAttr(resourceName, "name", "default"),
						resource.TestCheckResourceAttr(resourceName, "path", "default"),
					),

					// Fields related to this type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "blob_count", "0"),
						resource.TestCheckResourceAttrSet(resourceName, "total_size_in_bytes"),
						resource.TestCheckResourceAttrSet(resourceName, "available_space_in_bytes"),
					),
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
