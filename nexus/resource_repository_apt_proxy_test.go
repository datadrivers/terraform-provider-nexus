package nexus

import (
	"strconv"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testAccResourceRepositoryAptProxy() nexus.Repository {
	repo := testAccResourceRepositoryProxy(nexus.RepositoryFormatApt)
	repo.RepositoryApt = &nexus.RepositoryApt{
		Distribution: "bionic",
		Flat:         true,
	}
	useTrustStore := true
	repo.RepositoryProxy.RemoteURL = "https://remote.repository.com"
	repo.RepositoryHTTPClient.Connection = &nexus.RepositoryHTTPClientConnection{
		UseTrustStore: &useTrustStore,
	}
	return repo
}

func TestAccResourceRepositoryAptProxy(t *testing.T) {
	repo := testAccResourceRepositoryAptProxy()
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
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "bower.#", "0"),
						resource.TestCheckResourceAttr(resName, "maven.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker_proxy.#", "0"),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "apt.#", "1"),
						resource.TestCheckResourceAttr(resName, "apt.0.distribution", repo.RepositoryApt.Distribution),
						resource.TestCheckResourceAttr(resName, "apt.0.flat", strconv.FormatBool(repo.RepositoryApt.Flat)),
					),
				),
			},
			{
				ResourceName:            resName,
				ImportStateId:           repo.Name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"http_client.0.authentication.0.password"},
				// TODO: add check for storage
				// TODO: add check for apt.distribution
				// TODO: add tests for readonly repository
			},
		},
	})
}
