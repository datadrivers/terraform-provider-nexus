package nexus

import (
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testAccResourceRepositoryAptHosted() nexus.Repository {
	repo := testAccResourceRepositoryHosted(nexus.RepositoryFormatApt)
	repo.RepositoryApt = &nexus.RepositoryApt{
		Distribution: "bionic",
	}
	repo.RepositoryAptSigning = &nexus.RepositoryAptSigning{
		Keypair:    acctest.RandString(10),
		Passphrase: acctest.RandString(10),
	}
	return repo
}

func TestAccResourceRepositoryAptHosted(t *testing.T) {
	repo := testAccResourceRepositoryAptHosted()
	resName := testAccResourceRepositoryName(repo)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
						resource.TestCheckResourceAttr(resName, "apt.0.distribution", repo.RepositoryApt.Distribution),
						resource.TestCheckResourceAttr(resName, "apt_signing.#", "1"),
						resource.TestCheckResourceAttr(resName, "apt_signing.0.keypair", repo.RepositoryAptSigning.Keypair),
						resource.TestCheckResourceAttr(resName, "apt_signing.0.passphrase", repo.RepositoryAptSigning.Passphrase),
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
