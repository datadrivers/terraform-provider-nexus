package security_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/dre2004/go-nexus-client/nexus3/schema/security"
	"github.com/dre2004/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccResourceSecurityUser() security.User {
	return security.User{
		UserID:       fmt.Sprintf("user-test-%s", acctest.RandString(10)),
		FirstName:    fmt.Sprintf("user-firstname-%s", acctest.RandString(10)),
		LastName:     fmt.Sprintf("user-lastname-%s", acctest.RandString(10)),
		EmailAddress: fmt.Sprintf("user-email-%s@example.com", acctest.RandString(10)),
		Status:       "active",
		Password:     acctest.RandString(16),
		Roles:        []string{"nx-admin"},
	}
}

func TestAccResourceSecurityUser(t *testing.T) {
	resName := "nexus_security_user.acceptance"

	user := testAccResourceSecurityUser()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityUserConfig(user),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr(resName, "id", user.UserID),
					resource.TestCheckResourceAttr(resName, "userid", user.UserID),
					resource.TestCheckResourceAttr(resName, "firstname", user.FirstName),
					resource.TestCheckResourceAttr(resName, "lastname", user.LastName),
					resource.TestCheckResourceAttr(resName, "password", user.Password),
					resource.TestCheckResourceAttr(resName, "email", user.EmailAddress),
					resource.TestCheckResourceAttr(resName, "status", user.Status),
					resource.TestCheckResourceAttr(resName, "roles.#", strconv.Itoa(len(user.Roles))),
					// FIXME: (BUG) Incorrect roles state representation.
					// For some reasons, 1st element in array is not stored as roles.0, but instead it's stored
					// as roles.3360874991 where 3360874991 is a "random" number.
					// This number changes from test run to test run.
					// It may be a pointer to int instead of int itself, but it's not clear and requires additional research.
					// resource.TestCheckResourceAttr(resName, "roles.3360874991", "nx-admin"),
				),
			},
			{
				ResourceName:      resName,
				ImportStateId:     user.UserID,
				ImportState:       true,
				ImportStateVerify: true,
				// Password is not returned
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccResourceSecurityUserConfig(user security.User) string {
	return fmt.Sprintf(`
resource "nexus_security_user" "acceptance" {
	userid    = "%s"
	firstname = "%s"
	lastname  = "%s"
	email     = "%s"
	password  = "%s"
	status    = "%s"
	roles     = ["%s"]
}
`, user.UserID, user.FirstName, user.LastName, user.EmailAddress, user.Password, user.Status, strings.Join(user.Roles, "\", \""))
}
