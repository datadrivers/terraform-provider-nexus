package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceRealmsActive(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRealmsActive(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nexus_realms_active.acceptance", "realms.#", "2"),
					resource.TestCheckResourceAttr("data.nexus_realms_active.acceptance", "realms.0", "NexusAuthenticatingRealm"),
					resource.TestCheckResourceAttr("data.nexus_realms_active.acceptance", "realms.1", "NexusAuthorizingRealm"),
				),
			},
		},
	})

}

func testAccDataSourceRealmsActive() string {
	return fmt.Sprintf(`data "nexus_realms_active" "acceptance" {}`)
}
