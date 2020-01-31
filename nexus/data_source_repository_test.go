package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceRepository(t *testing.T) {
	repoName := "maven-central"
	resourceName := "data.nexus_repository.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRepository(repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "format"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "online"),
					// resource.TestCheckResourceAttrSet(resourceName, "storage"),
				),
			},
		},
	})
}

func testAccDataSourceRepository(name string) string {
	return fmt.Sprintf(`
data "nexus_repository" "acceptance" {
	name   = "%s"
}`, name)
}
