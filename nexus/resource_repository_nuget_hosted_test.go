package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRepositoryNugetHosted(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{

		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourceNugetHosted(repoName),
				Check: resource.ComposeTestCheckFunc(
					// Base and common repo props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "id", repoName),
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "name", repoName),
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "format", "nuget"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "type", "hosted"),
					),
					// Common fields
					// Online
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "online", "true"),
						// Storage
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "storage.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "storage.0.blob_store_name", "default"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "storage.0.strict_content_type_validation", "true"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "storage.0.write_policy", "ALLOW"),
					),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "maven.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "apt.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "apt_signing.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "bower.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "docker.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "docker_proxy.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "http_client.#", "0"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "group.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "negative_cache.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.nuget_hosted", "proxy.#", "0"),
					),
					// Fields related to this format and type
					// Format
					// - No special fields
					// Type
					resource.ComposeAggregateTestCheckFunc(
					// No specific fields
					),
				),
			},
			{
				ResourceName:      "nexus_repository.nuget_hosted",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func createTfStmtForResourceNugetHosted(name string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "nuget_hosted" {
	name   = "%s"
	format = "nuget"
	type   = "hosted"

	storage {
		write_policy = "ALLOW"
	}
}`, name)
}
