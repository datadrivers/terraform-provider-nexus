package nexus

import (
	"fmt"
	"strconv"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccResourceBlobstoreFile(t *testing.T) {
	resName := "nexus_blobstore.acceptance"

	bs := nexus.Blobstore{
		Name: fmt.Sprintf("test-blobstore-%d", acctest.RandIntRange(0, 99)),
		Type: nexus.BlobstoreTypeFile,
		Path: "/nexus-data/acceptance",
		BlobstoreSoftQuota: &nexus.BlobstoreSoftQuota{
			Limit: acctest.RandIntRange(100, 300) * 1000000,
			Type:  "spaceRemainingQuota",
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBlobstoreFileConfig(bs),
				Check: resource.ComposeTestCheckFunc(
					// Base and common resource props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "id", bs.Name),
						resource.TestCheckResourceAttr(resName, "name", bs.Name),
						resource.TestCheckResourceAttr(resName, "path", bs.Path),
						resource.TestCheckResourceAttr(resName, "type", bs.Type),
					),
					// Common fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "soft_quota.#", "1"),
						resource.TestCheckResourceAttr(resName, "soft_quota.0.limit", strconv.Itoa(bs.BlobstoreSoftQuota.Limit)),
						resource.TestCheckResourceAttr(resName, "soft_quota.0.type", bs.BlobstoreSoftQuota.Type),
					),
					// No fields related to other types
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "bucket_configuration.#", "0"),
					),

					// Fields related to this type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "blob_count", "0"),          // empty
						resource.TestCheckResourceAttr(resName, "total_size_in_bytes", "0"), // empty
						// FIXME: The value is unavailable, but should be
						// TODO: check that value is non-zero
						// resource.TestCheckResourceAttrSet(resName, "available_space_in_bytes"),
					),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateId:           bs.Name,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"available_space_in_bytes"},
			},
		},
	})
}

func testAccResourceBlobstoreFileConfig(bs nexus.Blobstore) string {
	return fmt.Sprintf(`
resource "nexus_blobstore" "acceptance" {
	name = "%s"
	path = "%s"
	type = "%s"

	soft_quota {
		limit = %d
		type  = "%s"
	}
}`, bs.Name, bs.Path, bs.Type, bs.BlobstoreSoftQuota.Limit, bs.BlobstoreSoftQuota.Type)
}
