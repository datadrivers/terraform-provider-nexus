package security_test

import (
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSecurityAnonymous(t *testing.T) {
	dataSourceName := "data.nexus_security_anonymous.acceptance"

	anonym := security.AnonymousAccessSettings{
		Enabled:   true,
		UserID:    acctest.RandString(20),
		RealmName: "NexusAuthenticatingRealm",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityAnonymousConfig(anonym) + testAccDataSourceSecurityAnonymousConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "enabled", strconv.FormatBool(anonym.Enabled)),
					resource.TestCheckResourceAttr(dataSourceName, "user_id", anonym.UserID),
					resource.TestCheckResourceAttr(dataSourceName, "realm_name", anonym.RealmName),
				),
			},
		},
	})
}

func testAccDataSourceSecurityAnonymousConfig() string {
	return `
data "nexus_security_anonymous" "acceptance" {
	depends_on = [
		nexus_security_anonymous.acceptance
	]
}
`
}
