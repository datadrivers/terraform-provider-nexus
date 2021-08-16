package nexus

import (
	"fmt"
	"strconv"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccResourceAnonymous(t *testing.T) {
	resName := "nexus_anonymous.acceptance"

	anonym := nexus.AnonymousConfig{
		Enabled:   true,
		UserID:    acctest.RandString(20),
		RealmName: acctest.RandString(20),
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
			{
				ResourceName:      resName,
				ImportStateId:     "anonymous",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceAnonymousConfig(anonym nexus.AnonymousConfig) string {
	return fmt.Sprintf(`
resource "nexus_anonymous" "acceptance" {
	enabled = "%t"
	user_id   = "%s"
	realm_name = "%s"
}
`, anonym.Enabled, anonym.UserID, anonym.RealmName)
}
