package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceRepositoryMavenProxy(t *testing.T) {
	repoName := "maven-central"
	resourceName := "data.nexus_repository.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRepository(repoName),
				Check: resource.ComposeTestCheckFunc(
					// Base and common repo props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", repoName),
						resource.TestCheckResourceAttr(resourceName, "name", repoName),
						resource.TestCheckResourceAttr(resourceName, "format", "maven"),
						resource.TestCheckResourceAttr(resourceName, "type", "proxy"),
					),
					// Common fields
					// Online
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "online", "true"),
						// Storage
						resource.TestCheckResourceAttr(resourceName, "storage.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "storage.0.blob_store_name", "default"),
						resource.TestCheckResourceAttr(resourceName, "storage.0.strict_content_type_validation", "false"),
						resource.TestCheckResourceAttr(resourceName, "storage.0.write_policy", "ALLOW"),
					),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "apt.#", "0"),
						resource.TestCheckResourceAttr(resourceName, "apt_signing.#", "0"),
						resource.TestCheckResourceAttr(resourceName, "bower.#", "0"),
						resource.TestCheckResourceAttr(resourceName, "docker.#", "0"),
						resource.TestCheckResourceAttr(resourceName, "docker_proxy.#", "0"),
						resource.TestCheckResourceAttr(resourceName, "http_client.#", "0"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "group.#", "0"),
					),
					// Fields related to this format and type
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "maven.#", "1"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "proxy.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "proxy.0.content_max_age", "-1"),
						resource.TestCheckResourceAttr(resourceName, "proxy.0.metadata_max_age", "1440"),
						resource.TestCheckResourceAttr(resourceName, "proxy.0.remote_url", "https://repo1.maven.org/maven/"),
						resource.TestCheckResourceAttr(resourceName, "negative_cache.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "negative_cache.0.enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "negative_cache.0.ttl", "1440"),
					),
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
