package deprecated_test

import (
	"testing"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRepositoryMavenProxy(t *testing.T) {
	repoName := "maven-central"
	resourceName := "data.nexus_repository.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRepositoryConfig(repoName),
				Check: resource.ComposeTestCheckFunc(
					// Base and common repo props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", repoName),
						resource.TestCheckResourceAttr(resourceName, "name", repoName),
						resource.TestCheckResourceAttr(resourceName, "format", "maven2"),
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
					),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "apt.#", "0"),
						resource.TestCheckResourceAttr(resourceName, "apt_signing.#", "0"),
						resource.TestCheckResourceAttr(resourceName, "bower.#", "0"),
						resource.TestCheckResourceAttr(resourceName, "docker.#", "0"),
						resource.TestCheckResourceAttr(resourceName, "docker_proxy.#", "0"),
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
						resource.TestCheckResourceAttr(resourceName, "http_client.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "http_client.0.authentication.#", "0"),
						resource.TestCheckResourceAttr(resourceName, "http_client.0.auto_block", "false"),
						resource.TestCheckResourceAttr(resourceName, "http_client.0.blocked", "false"),
						resource.TestCheckResourceAttr(resourceName, "proxy.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "proxy.0.content_max_age", "-1"),
						resource.TestCheckResourceAttr(resourceName, "proxy.0.metadata_max_age", "1440"),
						resource.TestCheckResourceAttr(resourceName, "proxy.0.remote_url", "https://repo1.maven.org/maven2/"),
						resource.TestCheckResourceAttr(resourceName, "negative_cache.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "negative_cache.0.enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "negative_cache.0.ttl", "1440"),
					),
				),
			},
		},
	})
}
