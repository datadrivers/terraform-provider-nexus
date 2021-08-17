package nexus

import (
	"strconv"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceSecurityUserToken(t *testing.T) {
	dataSourceName := "data.nexus_security_user_token.acceptance"

	token := nexus.UserTokenConfiguration{
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
