package security_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/dre2004/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSecurityLDAPOrder(t *testing.T) {
	resName := "nexus_security_ldap_order.acceptance"
	ldap := testAccResourceSecurityLDAP()
	log.Println(testAccResourceSecurityLDAPConfig(ldap) + testAccResourceSecurityLDAPOrder([]string{fmt.Sprintf("nexus_security_ldap.%s.name", ldap.Name)}))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityLDAPConfig(ldap) + testAccResourceSecurityLDAPOrder([]string{fmt.Sprintf("nexus_security_ldap.%s.name", ldap.Name)}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "order.#", "1"),
					resource.TestCheckResourceAttr(resName, "order.0", ldap.Name),
				),
			},
		},
	},
	)
}

func testAccResourceSecurityLDAPOrder(order []string) string {
	return fmt.Sprintf(`
resource "nexus_security_ldap_order" "acceptance" {
	order = [%s]
}
`, strings.Join(order, ", "))
}
