package nexus

import (
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSecurityUserToken(t *testing.T) {
	dataSourceName := "data.nexus_security_user_token.acceptance"

	token := security.UserTokenConfiguration{
		Enabled:        true,
		ProtectContent: false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityUserTokenConfig(token),
				Check:  nil,
			},
			{
				Config: testAccResourceSecurityUserTokenConfig(token) + testAccDataSourceSecurityUserTokenConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "enabled", strconv.FormatBool(token.Enabled)),
					resource.TestCheckResourceAttr(dataSourceName, "protect_content", strconv.FormatBool(token.ProtectContent)),
				),
			},
		},
	})
}

func testAccDataSourceSecurityUserTokenConfig() string {
	return `
data "nexus_security_user_token" "acceptance" {}
`
}
