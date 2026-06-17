package security_test

import (
	"testing"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSecurityOIDC(t *testing.T) {
	if tools.GetEnv("SKIP_PRO_TESTS", "false") == "true" {
		t.Skip("Skipping Nexus Pro tests")
	}

	cfg := testAccResourceSecurityOIDC()
	dsName := "data.nexus_security_oidc.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityOIDCConfig(cfg),
				Check:  nil,
			},
			{
				Config: testAccResourceSecurityOIDCConfig(cfg) + testAccDataSourceSecurityOIDCConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "client_id", cfg.ClientID),
					resource.TestCheckResourceAttr(dsName, "authorization_url", cfg.IdpAuthorizationURL),
					resource.TestCheckResourceAttr(dsName, "token_url", cfg.IdpTokenURL),
					resource.TestCheckResourceAttr(dsName, "jwks_url", cfg.IdpJwksURL),
					resource.TestCheckResourceAttr(dsName, "jws_algorithm", cfg.IdpJwsAlgorithm),
					resource.TestCheckResourceAttr(dsName, "username_claim", cfg.UsernameClaim),
					resource.TestCheckResourceAttr(dsName, "groups_claim", cfg.GroupsClaim),
				),
			},
		},
	})
}

func testAccDataSourceSecurityOIDCConfig() string {
	return `
data "nexus_security_oidc" "acceptance" {
}
`
}
