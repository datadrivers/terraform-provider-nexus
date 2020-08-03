package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceSecurityRealms(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecurityRealmsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nexus_security_realms.acceptance", "active.#", "2"),
					resource.TestCheckResourceAttrSet("data.nexus_security_realms.acceptance", "available.#"),
				),
			},
		},
	})

}

func testAccDataSourceSecurityRealmsConfig() string {
	return fmt.Sprintf(`
data "nexus_security_realms" "acceptance" {}
`)
}
