package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestDataSourceBlobstore(t *testing.T) {
	bsName := "default"
	resourceName := "data.nexus_blobstore.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBlobstore(bsName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "available_space_in_bytes"),
					resource.TestCheckResourceAttrSet(resourceName, "blob_count"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "path"),
					resource.TestCheckResourceAttrSet(resourceName, "total_size_in_bytes"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
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
