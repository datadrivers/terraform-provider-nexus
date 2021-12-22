package nexus

import (
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testAccResourceRepositoryHelmProxy() repository.LegacyRepository {
	repo := testAccResourceRepositoryProxy(repository.RepositoryFormatHelm)
	repo.RepositoryProxy.RemoteURL = "https://kubernetes-charts.storage.googleapis.com/"
	return repo
}

func TestAccResourceRepositoryHelmProxy(t *testing.T) {
	repo := testAccResourceRepositoryHelmProxy()
	resName := testAccResourceRepositoryName(repo)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryConfig(repo),
				Check: resource.ComposeTestCheckFunc(
					resourceRepositoryTestCheckFunc(repo),
					resourceRepositoryTypeProxyTestCheckFunc(repo),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "maven.#", "0"),
						resource.TestCheckResourceAttr(resName, "apt.#", "0"),
						resource.TestCheckResourceAttr(resName, "apt_signing.#", "0"),
						resource.TestCheckResourceAttr(resName, "bower.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker_proxy.#", "0"),
					),
				),
			},
			{
				ResourceName:            resName,
				ImportStateId:           repo.Name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"http_client.0.authentication.0.password"},
			},
		},
	})
}
