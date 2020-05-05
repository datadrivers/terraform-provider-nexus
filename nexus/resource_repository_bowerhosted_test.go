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
					// Base and common repo props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "id", repoName),
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "name", repoName),
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "format", "bower"),
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "type", "hosted"),
					),
					// Common fields
					// Online
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "online", "true"),
						// Storage
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "storage.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "storage.0.blob_store_name", "default"),
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "storage.0.strict_content_type_validation", "true"),
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "storage.0.write_policy", "ALLOW"),
					),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "maven.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "apt.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "apt_signing.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "docker.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "docker_proxy.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "http_client.#", "0"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "group.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "negative_cache.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "proxy.#", "0"),
					),
					// Fields related to this format and type
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "bower.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.bower_hosted", "bower.0.rewrite_package_urls", "true"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
					// No specific fields
					),
				),
			},
			{
				ResourceName:      "nexus_repository.bower_hosted",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
				// TODO: verify bower configuration, bower attribute is not returned by API currently
				ImportStateVerifyIgnore: []string{"bower"},
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
