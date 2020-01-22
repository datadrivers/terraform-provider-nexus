package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceUser(t *testing.T) {
	userID := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUser(userID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceUser("data.nexus_user.test", "nexus_user.test"),
				),
			},
		},
	})
}

func testAccCheckUser(userID string) string {
	return fmt.Sprintf(`
resource "nexus_user" "test" {
  userid    = "%s"
  firstname = "terraform-test"
  lastname  = "terraform-test"
  email     = "terraform-test@example.com"  
  password  = "test123"
  status    = "active"
  roles     = ["nx-admin"]
}

data "nexus_user" "test" {
  userid = nexus_user.test.userid
}
`, userID)
}

func testAccCheckDataSourceUser(data_source_name string, resource_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return nil
	}
}
