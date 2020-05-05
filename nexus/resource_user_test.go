package nexus

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccUser_update(t *testing.T) {
	// t.Parallel()

	userID := fmt.Sprintf("user-test-%s", acctest.RandString(10))
	userFirstname := fmt.Sprintf("user-firstname-%s", acctest.RandString(10))
	userLastname := fmt.Sprintf("user-lastname-%s", acctest.RandString(10))
	userEmail := fmt.Sprintf("user-email-%s@example.com", acctest.RandString(10))
	userPassword := acctest.RandString(16)
	userRoles := []string{"nx-admin"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserResource(userID, userFirstname, userLastname, userEmail, userPassword, "active", userRoles),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("nexus_user.acceptance", "id", userID),
					resource.TestCheckResourceAttr("nexus_user.acceptance", "userid", userID),
					resource.TestCheckResourceAttr("nexus_user.acceptance", "firstname", userFirstname),
					resource.TestCheckResourceAttr("nexus_user.acceptance", "lastname", userLastname),
					resource.TestCheckResourceAttr("nexus_user.acceptance", "password", userPassword),
					resource.TestCheckResourceAttr("nexus_user.acceptance", "email", userEmail),
					resource.TestCheckResourceAttr("nexus_user.acceptance", "status", "active"),
					resource.TestCheckResourceAttr("nexus_user.acceptance", "roles.#", "1"),
					// FIXME: (BUG) Incorrect roles state representation.
					// For some reasons, 1st element in array is not stored as roles.0, but instead it's stored
					// as roles.3360874991 where 3360874991 is a "random" number.
					// This number changes from test run to test run.
					// It may be a pointer to int instead of int itself, but it's not clear and requires additional research.
					// resource.TestCheckResourceAttr("nexus_user.acceptance", "roles.3360874991", "nx-admin"),
				),
			},
			{
				ResourceName:      "nexus_user.acceptance",
				ImportStateId:     userID,
				ImportState:       true,
				ImportStateVerify: true,
				// Password is not returned
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccUserResource(userID string, firstname string, lastname string, email string, password string, status string, roles []string) string {
	return fmt.Sprintf(`
resource "nexus_user" "acceptance" {
	userid    = "%s"
	firstname = "%s"
	lastname  = "%s"
	email     = "%s"
	password  = "%s"
	status    = "%s"
	roles     = ["%s"]
}
`, userID, firstname, lastname, email, password, status, strings.Join(roles, "\", \""))
}
