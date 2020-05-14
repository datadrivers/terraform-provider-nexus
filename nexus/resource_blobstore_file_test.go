package nexus

import (
	"fmt"
	"strconv"
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
				Config: testAccBlobstoreResourceFile(bsName, bsType, bsPath, quotaLimit, quotaType),
				Check: resource.ComposeTestCheckFunc(
					// Base and common resource props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_blobstore.acceptance", "id", bsName),
						resource.TestCheckResourceAttr("nexus_blobstore.acceptance", "name", bsName),
						resource.TestCheckResourceAttr("nexus_blobstore.acceptance", "path", bsPath),
						resource.TestCheckResourceAttr("nexus_blobstore.acceptance", "type", bsType),
					),
					// Common fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_blobstore.acceptance", "soft_quota.#", "1"),
						resource.TestCheckResourceAttr("nexus_blobstore.acceptance", "soft_quota.0.limit", strconv.Itoa(quotaLimit)),
						resource.TestCheckResourceAttr("nexus_blobstore.acceptance", "soft_quota.0.type", quotaType),
					),
					// No fields related to other types
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_blobstore.acceptance", "bucket_configuration.#", "0"),
					),

					// Fields related to this type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_blobstore.acceptance", "blob_count", "0"),          // empty
						resource.TestCheckResourceAttr("nexus_blobstore.acceptance", "total_size_in_bytes", "0"), // empty
						// FIXME: The value is unavailable, but should be
						// TODO: check that value is non-zero
						// resource.TestCheckResourceAttrSet("nexus_blobstore.acceptance", "available_space_in_bytes"),
					),
				),
			},
			{
				ResourceName:            "nexus_blobstore.acceptance",
				ImportState:             true,
				ImportStateId:           bsName,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"available_space_in_bytes"},
			},
		},
	})
}

func testAccBlobstoreResourceFile(name string, bsType string, path string, quotaLimit int, quotaType string) string {
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
