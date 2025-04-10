package cleanup_policy_test

//
//import (
//	"fmt"
//	"strconv"
//	_ "strings"
//	"testing"
//
//	nexusSchema "github.com/datadrivers/go-nexus-client/nexus3/schema/cleanuppolicies"
//	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
//)
//
//func TestAccResourceCleanupPolicy_basic(t *testing.T) {
//	resourceName := "nexus_cleanup_policy.test"
//
//	policy := nexusSchema.CleanupPolicy{
//		Name:                    acctest.RandString(20),
//		Notes:                   acctest.RandString(50),
//		CriteriaLastBlobUpdated: acctest.RandInt(),
//		CriteriaLastDownloaded:  acctest.RandInt(),
//		CriteriaReleaseType:     "RELEASES",
//		CriteriaAssetRegex:      ".*",
//		Retain:                  acctest.RandInt(),
//		Format:                  "maven2",
//	}
//
//	resource.Test(t, resource.TestCase{
//		Providers: acceptance.TestAccProviders,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccResourceCleanupPolicyCreateConfig(policy),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "name", policy.Name),
//					resource.TestCheckResourceAttr(resourceName, "notes", policy.Notes),
//					resource.TestCheckResourceAttr(resourceName, "criteria_last_blob_updated", strconv.Itoa(policy.CriteriaLastBlobUpdated)),
//					resource.TestCheckResourceAttr(resourceName, "criteria_last_downloaded", strconv.Itoa(policy.CriteriaLastDownloaded)),
//					resource.TestCheckResourceAttr(resourceName, "criteria_release_type", policy.CriteriaReleaseType),
//					resource.TestCheckResourceAttr(resourceName, "criteria_asset_regex", policy.CriteriaAssetRegex),
//					resource.TestCheckResourceAttr(resourceName, "retain", strconv.Itoa(policy.Retain)),
//					resource.TestCheckResourceAttr(resourceName, "format", policy.Format),
//				),
//			},
//		},
//	})
//}
//
//func testAccResourceCleanupPolicyCreateConfig(policy nexusSchema.CleanupPolicy) string {
//	return fmt.Sprintf(`
//resource "nexus_cleanup_policy" "test" {
//  name = "%s"
//  notes = "%s"
//  criteria_last_blob_updated = %d
//  criteria_last_downloaded = %d
//  criteria_release_type = "%s"
//  criteria_asset_regex = "%s"
//  retain = %d
//  format = "%s"
//}
//`, policy.Name, policy.Notes, policy.CriteriaLastBlobUpdated, policy.CriteriaLastDownloaded, policy.CriteriaReleaseType, policy.CriteriaAssetRegex, policy.Retain, policy.Format)
//}
