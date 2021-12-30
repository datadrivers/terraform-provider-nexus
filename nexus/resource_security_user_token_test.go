package nexus

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSecurityUserToken(t *testing.T) {
	resName := "nexus_security_user_token.acceptance"

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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "enabled", strconv.FormatBool(token.Enabled)),
					resource.TestCheckResourceAttr(resName, "protect_content", strconv.FormatBool(token.ProtectContent)),
				),
			},
		},
	})
}

func testAccResourceSecurityUserTokenConfig(token security.UserTokenConfiguration) string {
	return fmt.Sprintf(`
resource "nexus_security_user_token" "acceptance" {
	enabled         = %t
	protect_content = %t
}
`, token.Enabled, token.ProtectContent)
}
