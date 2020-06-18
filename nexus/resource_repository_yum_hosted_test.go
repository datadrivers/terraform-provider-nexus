package nexus

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRepositoryYumHosted(t *testing.T) {
	t.Parallel()

	repoName := fmt.Sprintf("test-repo-%s", acctest.RandString(10))
	repoYumDeployPolicy := "PERMISSIVE"
	repoYumRepodataDepth := 0

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourceYumHosted(repoName, repoYumDeployPolicy, repoYumRepodataDepth),
				Check: resource.ComposeTestCheckFunc(
					// Base and common repo props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "id", repoName),
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "name", repoName),
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "format", "yum"),
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "type", "hosted"),
					),
					// Common fields
					// Online
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "online", "true"),
						// Storage
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "storage.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "storage.0.blob_store_name", "default"),
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "storage.0.strict_content_type_validation", "true"),
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "storage.0.write_policy", "ALLOW"),
					),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "apt.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "bower.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "maven.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "docker.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "docker_proxy.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "http_client.#", "0"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "group.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "negative_cache.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "proxy.#", "0"),
					),
					// Fields related to this format and type
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "yum.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "yum.0.deploy_policy", repoYumDeployPolicy),
						resource.TestCheckResourceAttr("nexus_repository.yum_hosted", "yum.0.repodata_depth", strconv.Itoa(repoYumRepodataDepth)),
					),
				),
			},
			{
				ResourceName:      "nexus_repository.yum_hosted",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func createTfStmtForResourceYumHosted(name string, yumDeployPolicy string, yumRepodataDepth int) string {
	return fmt.Sprintf(`
resource "nexus_repository" "yum_hosted" {
	name   = "%s"
	format = "yum"
	type   = "hosted"

    yum {
        deploy_policy = "%s"
		repodata_depth = %d
    }

	storage {
		write_policy = "ALLOW"
	}
}
`, name, yumDeployPolicy, yumRepodataDepth)
}
