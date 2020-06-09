package nexus

import (
	"fmt"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRepositoryRAWHosted(t *testing.T) {
	resName := "nexus_repository.raw_hosted"
	repoName := fmt.Sprintf("test-raw-hosted-repo-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{

		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourceRAWHosted(repoName),
				Check: resource.ComposeTestCheckFunc(
					// Base and common repo props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "id", repoName),
						resource.TestCheckResourceAttr(resName, "name", repoName),
						resource.TestCheckResourceAttr(resName, "format", nexus.RepositoryFormatRAW),
						resource.TestCheckResourceAttr(resName, "type", nexus.RepositoryTypeHosted),
					),
					// Common fields
					// Online
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "online", "true"),
						// Storage
						resource.TestCheckResourceAttr(resName, "storage.#", "1"),
						resource.TestCheckResourceAttr(resName, "storage.0.blob_store_name", "default"),
						resource.TestCheckResourceAttr(resName, "storage.0.strict_content_type_validation", "true"),
						resource.TestCheckResourceAttr(resName, "storage.0.write_policy", "ALLOW"),
					),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "maven.#", "0"),
						resource.TestCheckResourceAttr(resName, "apt.#", "0"),
						resource.TestCheckResourceAttr(resName, "apt_signing.#", "0"),
						resource.TestCheckResourceAttr(resName, "bower.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker_proxy.#", "0"),
						resource.TestCheckResourceAttr(resName, "http_client.#", "0"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "group.#", "0"),
						resource.TestCheckResourceAttr(resName, "negative_cache.#", "0"),
						resource.TestCheckResourceAttr(resName, "proxy.#", "0"),
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
				ResourceName:      resName,
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func createTfStmtForResourceRAWHosted(name string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "raw_hosted" {
	format = "%s"
	name   = "%s"
	online = true
	type   = "%s"

	storage {
		write_policy = "ALLOW"
	}
}`, nexus.RepositoryFormatRAW, name, nexus.RepositoryTypeHosted)
}
