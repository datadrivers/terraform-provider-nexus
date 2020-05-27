package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRepositoryNpmProxy(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{

		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourceNpmProxy(repoName),
				Check: resource.ComposeTestCheckFunc(
					// Base and common repo props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "id", repoName),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "name", repoName),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "format", "npm"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "type", "proxy"),
					),
					// Common fields
					// Online
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "online", "true"),
						// Storage
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "storage.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "storage.0.blob_store_name", "default"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "storage.0.strict_content_type_validation", "true"),
						// FIXME: (BUG) Write policy can not be set to ALLOW is not set
						// resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "storage.0.write_policy", "ALLOW"),
					),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "maven.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "apt.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "apt_signing.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "bower.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "docker.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "docker_proxy.#", "0"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "group.#", "0"),
					),
					// Fields related to this format and type
					// Format
					// - No special fields
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "proxy.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "negative_cache.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "http_client.#", "1"),
					),
				),
			},
			{
				ResourceName:      "nexus_repository.npm_proxy",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
				// FIXME: (BUG) There is an inconsistency in http_client
				ImportStateVerifyIgnore: []string{"http_client"},
			},
		},
	})
}

func TestAccRepositoryNpmProxyWithoutAuth(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{

		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourceNpmProxyWithoutAuth(repoName),
				Check: resource.ComposeTestCheckFunc(
					// Base and common repo props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "id", repoName),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "name", repoName),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "format", "npm"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "type", "proxy"),
					),
					// Common fields
					// Online
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "online", "true"),
						// Storage
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "storage.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "storage.0.blob_store_name", "default"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "storage.0.strict_content_type_validation", "true"),
						// FIXME: (BUG) Write policy can not be set to ALLOW is not set
						// resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "storage.0.write_policy", "ALLOW"),
					),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "maven.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "apt.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "apt_signing.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "bower.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "docker.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "docker_proxy.#", "0"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "group.#", "0"),
					),
					// Fields related to this format and type
					// Format
					// - No special fields
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "proxy.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "negative_cache.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.npm_proxy", "http_client.#", "1"),
					),
				),
			},
			{
				ResourceName:      "nexus_repository.npm_proxy",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
				// FIXME: (BUG) There is an inconsistency in http_client
				ImportStateVerifyIgnore: []string{"http_client"},
			},
		},
	})
}

func createTfStmtForResourceNpmProxy(name string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "npm_proxy" {
	name   = "%s"
	format = "npm"
	type   = "proxy"

	proxy {
		remote_url  = "https://npm.org"
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

func createTfStmtForResourceNpmProxyWithoutAuth(name string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "npm_proxy" {
	name   = "%s"
	format = "npm"
	type   = "proxy"

	proxy {
		remote_url  = "https://npm.org"
	}

	http_client {
	}

	negative_cache {
		enabled = true
		ttl     = 1440
	}

	storage {

	}
}`, name)
}
