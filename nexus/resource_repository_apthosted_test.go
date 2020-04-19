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
					resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "name", repoName),
					resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "format", "apt"),
					resource.TestCheckResourceAttr("nexus_repository.apt_hosted", "type", "hosted"),
				),
			},
			{
				ResourceName:      "nexus_repository.apt_hosted",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
				// apt_signing not returned by API
				ImportStateVerifyIgnore: []string{"apt_signing"},
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
