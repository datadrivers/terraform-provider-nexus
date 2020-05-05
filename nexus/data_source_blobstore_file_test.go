package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestDataSourceBlobstoreFile(t *testing.T) {
	bsName := "default"
	resourceName := "data.nexus_blobstore.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBlobstore(bsName),
				Check: resource.ComposeTestCheckFunc(
					// Base and common resource props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", "default"),
						resource.TestCheckResourceAttr(resourceName, "name", "default"),
						resource.TestCheckResourceAttr(resourceName, "path", "default"),
						resource.TestCheckResourceAttr(resourceName, "type", "File"),
					),

					// Fields related to this type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "blob_count", "0"),          // empty
						resource.TestCheckResourceAttr(resourceName, "total_size_in_bytes", "0"), // empty
						// TODO: check that value is non-zero
						resource.TestCheckResourceAttrSet(resourceName, "available_space_in_bytes"),
					),
				),
			},
		},
	})
}

func testAccDataSourceBlobstore(name string) string {
	return fmt.Sprintf(`
data "nexus_blobstore" "acceptance" {
	name = "%s"
}`, name)
}
