package nexus

import (
	"fmt"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccResourceBlobstoreFile(t *testing.T) {
	bsName := fmt.Sprintf("test-blobstore-%d", acctest.RandIntRange(0, 99))
	bsType := nexus.BlobstoreTypeFile
	bsPath := fmt.Sprintf("/nexus-data/%s", bsName)
	quotaLimit := acctest.RandIntRange(100, 300)
	quotaType := "spaceRemainingQuota"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBlobstoreResource(bsName, bsType, bsPath, quotaLimit, quotaType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nexus_blobstore.acceptance", "name", bsName),
					resource.TestCheckResourceAttr("nexus_blobstore.acceptance", "type", bsType),
				),
			},
		},
	})
}

func testAccBlobstoreResource(name string, bsType string, path string, quotaLimit int, quotaType string) string {
	return fmt.Sprintf(`
resource "nexus_blobstore" "acceptance" {
	name = "%s"
	path = "%s"
	type = "%s"

	soft_quota {
		limit = %d
		type  = "%s"
	}
}`, name, path, bsType, quotaLimit, quotaType)
}
