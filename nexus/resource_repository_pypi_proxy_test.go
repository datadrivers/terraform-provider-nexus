package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRepositoryPypiProxy(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{

		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourcePypiProxy(repoName),
				Check: resource.ComposeTestCheckFunc(
					// Base and common repo props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "id", repoName),
						resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "name", repoName),
						resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "format", "pypi"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "type", "proxy"),
					),
					// Common fields
					// Online
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "online", "true"),
						// Storage
						resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "storage.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "storage.0.blob_store_name", "default"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "storage.0.strict_content_type_validation", "true"),
						// FIXME: (BUG) Write policy can not be set to ALLOW is not set
						// resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "storage.0.write_policy", "ALLOW"),
					),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "maven.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "apt.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "apt_signing.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "bower.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "docker.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "docker_proxy.#", "0"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "group.#", "0"),
					),
					// Fields related to this format and type
					// Format
					// - No special fields
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "proxy.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "negative_cache.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_proxy", "http_client.#", "1"),
					),
				),
			},
			{
				ResourceName:      "nexus_repository.pypi_proxy",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
				// FIXME: (BUG) There is an inconsistency in http_client
				ImportStateVerifyIgnore: []string{"http_client"},
			},
		},
	})
}

func createTfStmtForResourcePypiProxy(name string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "pypi_proxy" {
	name   = "%s"
	format = "pypi"
	type   = "proxy"

	proxy {
		remote_url  = "https://pypi.org"
	}

	http_client {
		authentication {
			type = "username"
		}
	}

	negative_cache {
		enabled = true
		ttl     = 1440
	}

	storage {

	}
}`, name)
}
