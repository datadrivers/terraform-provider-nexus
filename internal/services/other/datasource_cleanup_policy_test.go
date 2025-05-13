package other_test

import (
	"strconv"
	"testing"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCleanupPolicy(t *testing.T) {
	resName := "data.nexus_cleanup_policy.acceptance"
	policy := newCleanupPolicy()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCleanupPolicyConfig(policy) + testAccDataSourceCleanupPolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", policy.Name),
					resource.TestCheckResourceAttr(resName, "format", string(policy.Format)),
					resource.TestCheckResourceAttr(resName, "criteria.0.last_blob_updated", strconv.Itoa(*policy.CriteriaLastBlobUpdated)),
					resource.TestCheckResourceAttr(resName, "criteria.0.last_downloaded", strconv.Itoa(*policy.CriteriaLastDownloaded)),
					resource.TestCheckResourceAttr(resName, "criteria.0.release_type", string(*policy.CriteriaReleaseType)),
					resource.TestCheckResourceAttr(resName, "criteria.0.asset_regex", *policy.CriteriaAssetRegex),
				),
			},
		},
	})
}

func testAccDataSourceCleanupPolicyConfig() string {
	return `
data "nexus_cleanup_policy" "acceptance" {
	name = nexus_cleanup_policy.acceptance.name
}
`
}
