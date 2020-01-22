package nexus

import (
	"fmt"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccUser_update(t *testing.T) {
	t.Parallel()

	var user nexus.User

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
				Config: testAccUserResource(userID, userFirstname, userLastname, userEmail, userPassword, "active"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserResourceExists("nexus_user.test", &user),
					testAccCheckUserValues(user, userID, userFirstname, userLastname, userEmail),
					resource.TestCheckResourceAttr("nexus_user.test", "firstname", userFirstname),
					resource.TestCheckResourceAttr("nexus_user.test", "lastname", userLastname),
					resource.TestCheckResourceAttr("nexus_user.test", "email", userEmail),
				),
			},
			// {
			// 	ResourceName:      "nexus_user.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccUserResource(userID string, firstname string, lastname string, email string, password string, status string) string {
	return fmt.Sprintf(`
resource "nexus_user" "test" {
	userid    = "%s"
	firstname = "%s"
	lastname  = "%s"
	email     = "%s"
	password  = "%s"
	status    = "%s"
	roles     = ["nx-admin"]
}
`, userID, firstname, lastname, email, password, status)
}

func testAccCheckUserResourceExists(name string, user *nexus.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		nexusClient := testAccProvider.Meta().(nexus.Client)
		result, err := nexusClient.UserRead(rs.Primary.ID)
		if err != nil {
			return err
		}
		*user = *result

		return nil // fmt.Errorf("User `%s' not found", rs.Primary.ID)
	}
}

func testAccCheckUserValues(user nexus.User, userID string, firstname string, lastname string, email string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if user.UserID != userID {
			return fmt.Errorf("bad userid, expected \"%s\", got: %v", userID, user.UserID)
		}

		if user.FirstName != firstname {
			return fmt.Errorf("bad firstname, expected \"%s\", got: %v", firstname, user.FirstName)
		}
		return nil
	}
}
