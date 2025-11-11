package security_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSecurityUser(t *testing.T) {
	resName := "data.nexus_security_user.acceptance"
	user := security.User{
		UserID:       fmt.Sprintf("user-test-%s", acctest.RandString(10)),
		FirstName:    fmt.Sprintf("user-firstname-%s", acctest.RandString(10)),
		LastName:     fmt.Sprintf("user-lastname-%s", acctest.RandString(10)),
		EmailAddress: fmt.Sprintf("user-email-%s@example.com", acctest.RandString(10)),
		Status:       "active",
		Password:     acctest.RandString(16),
		Roles:        []string{"nx-admin"},
		Source:       "default",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityUserConfig(user) + testAccDataSourceSecurityUserConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "id", user.UserID),
					resource.TestCheckResourceAttr(resName, "userid", user.UserID),
					resource.TestCheckResourceAttr(resName, "firstname", user.FirstName),
					resource.TestCheckResourceAttr(resName, "lastname", user.LastName),
					// Password is not returned by API
					// resource.TestCheckResourceAttr(resName, "password", user.Password),
					resource.TestCheckResourceAttr(resName, "email", user.EmailAddress),
					resource.TestCheckResourceAttr(resName, "status", user.Status),
					resource.TestCheckResourceAttr(resName, "roles.#", strconv.Itoa(len(user.Roles))),
					resource.TestCheckResourceAttr(resName, "source", user.Source),
				),
			},
		},
	})
}

func testAccDataSourceSecurityUserConfig() string {
	return `
data "nexus_security_user" "acceptance" {
	userid = nexus_security_user.acceptance.userid
}
`
}
