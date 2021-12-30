package nexus

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSecurityLDAP(t *testing.T) {
	resName := "data.nexus_security_ldap.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecurityLDAPConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "ldap.#"),
					resource.TestCheckResourceAttr(resName, "ldap.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceSecurityLDAPConfig() string {
	return `
data "nexus_security_ldap" "acceptance" {}
`
}
