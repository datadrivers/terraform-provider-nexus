package deprecated_test

import (
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccResourceRepositoryAptHosted() repository.LegacyRepository {
	repo := testAccResourceRepositoryHosted(repository.RepositoryFormatApt)
	repo.Apt = &repository.AptProxy{
		Distribution: "bionic",
	}
	keypair := acctest.RandString(10)
	passphrase := acctest.RandString(10)
	repo.AptSigning = &repository.AptSigning{
		Keypair:    &keypair,
		Passphrase: &passphrase,
	}
	return repo
}

func TestAccResourceRepositoryAptHosted(t *testing.T) {
	repo := testAccResourceRepositoryAptHosted()
	resName := testAccResourceRepositoryName(repo)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryConfig(repo),
				Check: resource.ComposeTestCheckFunc(
					resourceRepositoryTestCheckFunc(repo),
					resourceRepositoryTypeHostedTestCheckFunc(repo),
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "bower.#", "0"),
						resource.TestCheckResourceAttr(resName, "maven.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker_proxy.#", "0"),
					),
					// Type
					// Fields related to this format and type
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "apt.#", "1"),
						resource.TestCheckResourceAttr(resName, "apt.0.distribution", repo.Apt.Distribution),
						resource.TestCheckResourceAttr(resName, "apt_signing.#", "1"),
						resource.TestCheckResourceAttr(resName, "apt_signing.0.keypair", *repo.AptSigning.Keypair),
						resource.TestCheckResourceAttr(resName, "apt_signing.0.passphrase", *repo.AptSigning.Passphrase),
					),
				),
			},
			{
				ResourceName:      resName,
				ImportStateId:     repo.Name,
				ImportState:       true,
				ImportStateVerify: true,
				// TODO: verify apt_signing configuration, apt_signing not returned by API currently
				ImportStateVerifyIgnore: []string{"apt_signing", "cleanup.0.policy_names"},
			},
		},
	})
}
