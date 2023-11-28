package security_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/dre2004/go-nexus-client/nexus3/schema/security"
	"github.com/dre2004/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSecurityAnonymous(t *testing.T) {
	resName := "nexus_security_anonymous.acceptance"

	anonym := security.AnonymousAccessSettings{
		Enabled:   true,
		UserID:    "acctest",
		RealmName: "NexusAuthenticatingRealm",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityAnonymousConfig(anonym),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "enabled", strconv.FormatBool(anonym.Enabled)),
					resource.TestCheckResourceAttr(resName, "user_id", anonym.UserID),
					resource.TestCheckResourceAttr(resName, "realm_name", anonym.RealmName),
				),
			},
		},
	})
}

func testAccResourceSecurityAnonymousConfig(anonym security.AnonymousAccessSettings) string {
	return fmt.Sprintf(`
resource "nexus_security_anonymous" "acceptance" {
	enabled    = "%t"
	user_id    = "%s"
	realm_name = "%s"
}
`, anonym.Enabled, anonym.UserID, anonym.RealmName)
}
