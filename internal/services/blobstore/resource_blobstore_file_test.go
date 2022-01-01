package blobstore_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/blobstore"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceBlobstoreFile(t *testing.T) {
	resName := "nexus_blobstore_file.acceptance"

	bs := blobstore.File{
		Name: fmt.Sprintf("test-blobstore-%d", acctest.RandIntRange(0, 99)),
		Path: "/nexus-data/acceptance",
		SoftQuota: &blobstore.SoftQuota{
			Limit: int64(acctest.RandIntRange(100, 300) * 1000000),
			Type:  "spaceRemainingQuota",
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
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
					),
					// Common fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "soft_quota.#", "1"),
						resource.TestCheckResourceAttr(resName, "soft_quota.0.limit", strconv.FormatInt(bs.SoftQuota.Limit, 10)),
						resource.TestCheckResourceAttr(resName, "soft_quota.0.type", bs.SoftQuota.Type),
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

func testAccResourceBlobstoreFileConfig(bs blobstore.File) string {
	return fmt.Sprintf(`
resource "nexus_blobstore_file" "acceptance" {
	name = "%s"
	path = "%s"

	soft_quota {
		limit = %d
		type  = "%s"
	}
}`, bs.Name, bs.Path, bs.SoftQuota.Limit, bs.SoftQuota.Type)
}
