package cleanuppolicies_test

import (
	"github.com/datadrivers/go-nexus-client/nexus3/schema/cleanuppolicies"
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/pkg/tools"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCleanupPolicy(t *testing.T) {
	dataSourceName := "data.nexus_cleanup_policy.acceptance"

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
					resource.TestCheckResourceAttrPtr(dataSourceName, "notes", cp.Notes),
					resource.TestCheckResourceAttr(dataSourceName, "criteria_last_blob_updated", strconv.Itoa(*cp.CriteriaLastBlobUpdated)),
					resource.TestCheckResourceAttr(dataSourceName, "criteria_last_downloaded", strconv.Itoa(*cp.CriteriaLastDownloaded)),
					resource.TestCheckResourceAttrPtr(dataSourceName, "criteria_release_type", cp.CriteriaReleaseType),
					resource.TestCheckResourceAttrPtr(dataSourceName, "criteria_asset_regex", cp.CriteriaAssetRegex),
					resource.TestCheckResourceAttr(dataSourceName, "retain", strconv.Itoa(cp.Retain)),
					resource.TestCheckResourceAttr(dataSourceName, "name", cp.Name),
					resource.TestCheckResourceAttr(dataSourceName, "format", cp.Format),
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
