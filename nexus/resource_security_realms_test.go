package nexus

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSecurityRealms(t *testing.T) {
	resName := "nexus_security_realms.acceptance"
	realms := []string{"NexusAuthenticatingRealm", "NexusAuthorizingRealm"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityRealmsConfig(realms),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "active.#", strconv.Itoa(len(realms))),
					resource.TestCheckResourceAttr(resName, "active.0", realms[0]),
					resource.TestCheckResourceAttr(resName, "active.1", realms[1]),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceSecurityRealmsConfig(realms []string) string {
	return fmt.Sprintf(`
resource "nexus_security_realms" "acceptance" {
	active = ["%s"]
}`, strings.Join(realms, "\",\""))
}
