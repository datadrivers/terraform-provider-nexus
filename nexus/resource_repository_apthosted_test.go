package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRepositoryAptHosted(t *testing.T) {
	t.Parallel()

	repoName := fmt.Sprintf("test-repo-%s", acctest.RandString(10))
	repoAptDistribution := "bionic"
	repoAptSigningKeypair := acctest.RandString(10)
	repoAptSigningPassphrase := acctest.RandString(10)
	repoCleanupPolicyNames := []string{"weekly-cleanup"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourceAptHosted(repoName, repoAptDistribution, repoAptSigningKeypair, repoAptSigningPassphrase, repoCleanupPolicyNames),
				Check: resource.ComposeTestCheckFunc(
					// Base and common repo props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "id", repoName),
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "name", repoName),
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "format", "apt"),
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "type", "hosted"),
					),
					// Common fields
					// Online
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "online", "true"),
						// Storage
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "storage.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "storage.0.blob_store_name", "default"),
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "storage.0.strict_content_type_validation", "true"),
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "storage.0.write_policy", "ALLOW"),
					),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "bower.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "maven.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "docker.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "docker_proxy.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "http_client.#", "0"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "group.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "negative_cache.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "proxy.#", "0"),
					),
					// Fields related to this format and type
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "apt.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "apt.0.distribution", repoAptDistribution),
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "apt_signing.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "apt_signing.0.keypair", repoAptSigningKeypair),
						resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "apt_signing.0.passphrase", repoAptSigningPassphrase),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
					// No specific fields
					),
				),
			},
			{
				ResourceName:      "nexus_repository.apt_hosted",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
				// TODO: verify apt_signing configuration, apt_signing not returned by API currently
				ImportStateVerifyIgnore: []string{"apt_signing"},
				// TODO: add check for storage
				// TODO: add check for apt.distribution
				// TODO: add tests for readonly repository
			},
		},
	})
}

func createTfStmtForResourceAptHosted(name string, aptDistribution string, aptSigningKEypair string, aptSigningPassphrase string, cleanupPolicyNames []string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "apt_hosted" {
	name   = "%s"
	format = "apt"
	type   = "hosted"

	apt {
		distribution = "%s"
	}

	apt_signing {
		keypair    = "%s"
		passphrase = "%s"
	}

	storage {
		write_policy = "ALLOW"
	}
}
`, name, aptDistribution, aptSigningKEypair, aptSigningPassphrase)
}
