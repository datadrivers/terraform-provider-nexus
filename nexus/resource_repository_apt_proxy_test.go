package nexus

import (
	"fmt"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRepositoryAptProxy(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-%s", acctest.RandString(10))
	repoDistribution := "bionic"
	repoFlat := "true"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourceAptProxy(repoName, repoDistribution, repoFlat),
				Check: resource.ComposeTestCheckFunc(
					// Base and common repo props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "id", repoName),
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "name", repoName),
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "format", nexus.RepositoryFormatApt),
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "type", nexus.RepositoryTypeProxy),
					),
					// Common fields
					// Online
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "online", "true"),
						// Storage
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "storage.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "storage.0.blob_store_name", "default"),
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "storage.0.strict_content_type_validation", "true"),
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "storage.0.write_policy", "ALLOW"),
					),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "bower.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "maven.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "docker.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "docker_proxy.#", "0"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "group.#", "0"),
					),
					// Fields related to this format and type
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "apt.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "apt.0.distribution", repoDistribution),
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "apt.0.flat", repoFlat),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "negative_cache.#", "1"),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.apt_proxy", "proxy.#", "1"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
					// No specific fields
					),
				),
			},
			{
				ResourceName:      "nexus_repository.apt_proxy",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
				// TODO: add check for storage
				// TODO: add check for apt.distribution
				// TODO: add tests for readonly repository
			},
		},
	})
}

func createTfStmtForResourceAptProxy(name string, distribution string, flat string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "apt_proxy" {
	format = "%s"
	name   = "%s"
	online = true
	type   = "%s"

	apt {
		distribution = "%s"
		flat         = %s
	}

	http_client {
		auto_block = true
	}

	negative_cache {
		enabled = true
		ttl     = 1440
	}

	proxy {
		content_max_age  = "1440"
		metadata_max_age = "1440"
		remote_url       = "http://archive.ubuntu.com/ubuntu/"
	}

	storage {
		blob_store_name                = "default"
		strict_content_type_validation = true
		write_policy                   = "ALLOW"
	}
}
`, nexus.RepositoryFormatApt, name, nexus.RepositoryTypeProxy, distribution, flat)
}
