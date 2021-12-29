package nexus

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAnonymous(t *testing.T) {
	resName := "nexus_anonymous.acceptance"

	anonym := security.AnonymousAccessSettings{
		Enabled:   true,
		UserID:    "acctest",
		RealmName: "NexusAuthenticatingRealm",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAnonymousConfig(anonym),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "enabled", strconv.FormatBool(anonym.Enabled)),
					resource.TestCheckResourceAttr(resName, "user_id", anonym.UserID),
					resource.TestCheckResourceAttr(resName, "realm_name", anonym.RealmName),
				),
			},
		},
	})
}

func testAccResourceAnonymousConfig(anonym security.AnonymousAccessSettings) string {
	return fmt.Sprintf(`
resource "nexus_anonymous" "acceptance" {
	enabled    = "%t"
	user_id    = "%s"
	realm_name = "%s"
}
`, anonym.Enabled, anonym.UserID, anonym.RealmName)
}
