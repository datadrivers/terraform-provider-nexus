package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccUser_update(t *testing.T) {
	t.Parallel()

	userID := fmt.Sprintf("user-test-%s", acctest.RandString(10))
	userFirstname := fmt.Sprintf("user-firstname-%s", acctest.RandString(10))
	userLastname := fmt.Sprintf("user-lastname-%s", acctest.RandString(10))
	userEmail := fmt.Sprintf("user-email-%s@example.com", acctest.RandString(10))
	userPassword := acctest.RandString(16)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUser_basic(userID, userFirstname, userLastname, userEmail, userPassword),
			},
			{
				ResourceName:      "nexus_user.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccUser_basic(userID string, firstname string, lastname string, email string, password string) string {
	return fmt.Sprintf(`
resource "nexus_user" "foobar" {
	userid    = "%s"
	firstname = "%s"
	lastname  = "%s"
	email     = "%s"
	password  = "%s"
	status    = "active"
	roles     = ["nx-all"]
}
`, userID, firstname, lastname, email, password)
}
