package security_test

import (
	"testing"

	"github.com/dre2004/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSecurityLDAP(t *testing.T) {
	resName := "data.nexus_security_ldap.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
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
