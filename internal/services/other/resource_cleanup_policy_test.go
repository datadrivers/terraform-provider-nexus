package other_test

import (
	"fmt"
	"strconv"
	"testing"

	schema "github.com/datadrivers/go-nexus-client/nexus3/schema/cleanuppolicies"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func newCleanupPolicy() schema.CleanupPolicy {
	cr := schema.CriteriaReleaseTypeReleases
	policy := schema.CleanupPolicy{
		Name:                    acctest.RandString(10),
		Format:                  "maven2",
		CriteriaLastBlobUpdated: tools.GetIntPointer(acctest.RandIntRange(1, 9999)),
		CriteriaLastDownloaded:  tools.GetIntPointer(acctest.RandIntRange(1, 9999)),
		CriteriaAssetRegex:      tools.GetStringPointer(".*"),
		CriteriaReleaseType:     &cr,
	}

	return policy
}

func TestAccResourceCleanupPolicy(t *testing.T) {
	resName := "nexus_cleanup_policy.acceptance"
	policy := newCleanupPolicy()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCleanupPolicyConfig(policy),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "id", policy.Name),
					resource.TestCheckResourceAttr(resName, "name", policy.Name),
					resource.TestCheckResourceAttr(resName, "format", string(policy.Format)),
					resource.TestCheckResourceAttr(resName, "criteria.0.last_blob_updated", strconv.Itoa(*policy.CriteriaLastBlobUpdated)),
					resource.TestCheckResourceAttr(resName, "criteria.0.last_downloaded", strconv.Itoa(*policy.CriteriaLastDownloaded)),
					resource.TestCheckResourceAttr(resName, "criteria.0.release_type", string(*policy.CriteriaReleaseType)),
					resource.TestCheckResourceAttr(resName, "criteria.0.asset_regex", *policy.CriteriaAssetRegex),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateId:     policy.Name,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceCleanupPolicyConfig(policy schema.CleanupPolicy) string {
	return fmt.Sprintf(`
	resource "nexus_cleanup_policy" "acceptance" {
		name   = "%s"
		format = "%s"
		criteria {
			last_blob_updated = %d
			last_downloaded   = %d
			release_type      = "%s"
			asset_regex       = "%s"
		}
	  }
`,
		policy.Name,
		policy.Format,
		policy.CriteriaLastBlobUpdated,
		policy.CriteriaLastDownloaded,
		*policy.CriteriaReleaseType,
		*policy.CriteriaAssetRegex,
	)
}
