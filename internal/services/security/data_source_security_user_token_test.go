package security_test

import (
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSecurityUserToken(t *testing.T) {
	if tools.GetEnv("SKIP_PRO_TESTS", "false") == "true" {
		t.Skip("Skipping Nexus Pro tests")
	}

	dataSourceName := "data.nexus_security_user_token.acceptance"

	token := security.UserTokenConfiguration{
		Enabled:           true,
		ProtectContent:    false,
		ExpirationEnabled: true,
		ExpirationDays:    int(30),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
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
					resource.TestCheckResourceAttr(dataSourceName, "expiration_enabled", strconv.FormatBool(token.ExpirationEnabled)),
					resource.TestCheckResourceAttr(dataSourceName, "expiration_days", strconv.Itoa(token.ExpirationDays)),
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
