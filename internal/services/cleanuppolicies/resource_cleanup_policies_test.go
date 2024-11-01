package cleanuppolicies_test

import (
	"fmt"
	"strconv"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/pkg/tools"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/cleanuppolicies"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceCleanupPolicy(t *testing.T) {
	var cleanupPolicy cleanuppolicies.CleanupPolicy

	resName := "nexus_security_cleanup_policy.acceptance"
	cp := cleanuppolicies.CleanupPolicy{
		Notes:                   tools.GetStringPointer(acctest.RandString(25)),
		CriteriaLastBlobUpdated: tools.GetIntPointer(acctest.RandInt()),
		CriteriaLastDownloaded:  tools.GetIntPointer(acctest.RandInt()),
		CriteriaReleaseType:     tools.GetStringPointer(acctest.RandString(8)),
		CriteriaAssetRegex:      tools.GetStringPointer(acctest.RandString(15)),
		Retain:                  acctest.RandInt(),
		Name:                    acctest.RandString(10),
		Format:                  acctest.RandString(5),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCleanupPolicyConfig(cp) + testAccDataSourceCleanupPolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "notes", *cp.Notes),
					resource.TestCheckResourceAttr(resName, "criteria_last_blob_updated", strconv.Itoa(*cp.CriteriaLastBlobUpdated)),
					resource.TestCheckResourceAttr(resName, "criteria_last_downloaded", strconv.Itoa(*cp.CriteriaLastDownloaded)),
					resource.TestCheckResourceAttr(resName, "criteria_release_type", *cp.CriteriaReleaseType),
					resource.TestCheckResourceAttr(resName, "criteria_asset_regex", *cp.CriteriaAssetRegex),
					resource.TestCheckResourceAttr(resName, "retain", strconv.Itoa(cp.Retain)),
					resource.TestCheckResourceAttr(resName, "name", cp.Name),
					resource.TestCheckResourceAttr(resName, "format", cp.Format),
				),
			},
			{
				ResourceName:      resName,
				ImportStateId:     cleanupPolicy.Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceCleanupPolicyConfig(cp cleanuppolicies.CleanupPolicy) string {
	return fmt.Sprintf(`
resource "nexus_security_cleanup_policy" "acceptance" {
	notes = "%s"
	criteria_last_blob_updated = "%s"
	criteria_last_downloaded = "%s"
	criteria_release_type = "%s"
	criteria_asset_regex = "%s"
	retain = "%s"
	name = "%s"
	format = "%s"
}
`, *cp.Notes, strconv.Itoa(*cp.CriteriaLastBlobUpdated), strconv.Itoa(*cp.CriteriaLastDownloaded),
		*cp.CriteriaReleaseType, *cp.CriteriaAssetRegex, strconv.Itoa(cp.Retain), cp.Name, cp.Format)
}

func testAccCheckCleanupPolicyResourceExists(name string, cleanupPolicy *cleanuppolicies.CleanupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		client := acceptance.TestAccProvider.Meta().(*nexus.NexusClient)
		result, err := client.CleanupPolicy.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		*cleanupPolicy = *result

		return nil
	}
}
