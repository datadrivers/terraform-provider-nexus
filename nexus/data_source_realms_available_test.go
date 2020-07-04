package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceRealmsAvailable(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRealmsAvailable(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nexus_realms_available.acceptance", "realms.#"),
				),
			},
		},
	})
}

func testAccDataSourceRealmsAvailable() string {
	return fmt.Sprintf(`data "nexus_realms_available" "acceptance" {}`)
}
