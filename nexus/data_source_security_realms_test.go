package nexus

import (
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
					resource.TestCheckResourceAttrSet("data.nexus_security_realms.acceptance", "active.#"),
					resource.TestCheckResourceAttrSet("data.nexus_security_realms.acceptance", "available.#"),
				),
			},
		},
	})

}

func testAccDataSourceSecurityRealmsConfig() string {
	return `
data "nexus_security_realms" "acceptance" {}
`
}
