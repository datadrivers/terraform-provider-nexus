package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRepositoryNugetProxy(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-nuget-proxy-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourceNugetProxy(repoName),
				Check: resource.ComposeTestCheckFunc(
					// Base and common repo props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.nuget_proxy", "id", repoName),
						resource.TestCheckResourceAttr("nexus_repository.nuget_proxy", "name", repoName),
						resource.TestCheckResourceAttr("nexus_repository.nuget_proxy", "format", "nuget"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_proxy", "type", "proxy"),
					),
					// Common fields
					resource.ComposeAggregateTestCheckFunc(
						// Online
						resource.TestCheckResourceAttr("nexus_repository.nuget_proxy", "online", "true"),
						// Storage
						resource.TestCheckResourceAttr("nexus_repository.nuget_proxy", "storage.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_proxy", "storage.0.blob_store_name", "default"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_proxy", "storage.0.strict_content_type_validation", "true"),
					),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.nuget_proxy", "apt.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_proxy", "apt_signing.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_proxy", "bower.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_proxy", "docker.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_proxy", "docker_proxy.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_proxy", "maven.#", "0"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.nuget_proxy", "group.#", "0"),
					),
					// Fields related to this format and type
					// Format
					// - No special fields
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.nuget_proxy", "http_client.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_proxy", "negative_cache.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_proxy", "proxy.#", "1"),
					),
				),
			},
			{
				ResourceName:      "nexus_repository.nuget_proxy",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
				// FIXME: (BUG) There is an inconsistency in http_client
				ImportStateVerifyIgnore: []string{"http_client"},
			},
		},
	})
}

func createTfStmtForResourceNugetProxy(name string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "nuget_proxy" {
	name   = "%s"
	format = "nuget"
	type   = "proxy"
	online = true

	http_client {
		authentication {
			type = "username"
		}
	}

	negative_cache {
		enabled = true
		ttl     = 1440
	}

	nuget_proxy {
		query_cache_item_max_age = 1234
	}

	proxy {
		remote_url  = "https://www.nuget.org/api/v2/"
	}

	storage {
		write_policy = "ALLOW"
	}
}`, name)
}
