package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRepositoryBowerHosted(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-%s", acctest.RandString(10))
	bowerRewritePackageURLs := true

	resource.Test(t, resource.TestCase{

		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourceBowerHosted(repoName, bowerRewritePackageURLs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "name", repoName),
					resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "format", "bower"),
					resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "type", "hosted"),
				),
			},
			{
				ResourceName:      "nexus_repository.bower_hosted",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
				// TODO: verify bower configuration, bower attribute is not returned by API currently
				ImportStateVerifyIgnore: []string{"bower"},
				// TODO: add check for storage
				// TODO: add tests for readonly repository
			},
		},
	})
}

func createTfStmtForResourceBowerHosted(name string, rewritePackageURLs bool) string {
	return fmt.Sprintf(`
resource "nexus_repository" "bower_hosted" {
	name   = "%s"
	format = "bower"
	type   = "hosted"

	bower {
		rewrite_package_urls = %v
	}

	storage {
		write_policy = "ALLOW"
	}
}`, name, rewritePackageURLs)
}
