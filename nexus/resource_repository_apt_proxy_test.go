package nexus

import (
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccResourceRepositoryAptProxy() repository.LegacyRepository {
	repo := testAccResourceRepositoryProxy(repository.RepositoryFormatApt)
	repo.Apt = &repository.AptProxy{
		Distribution: "bionic",
		Flat:         true,
	}
	useTrustStore := true
	remoteURL := "https://remote.repository.com"
	repo.Proxy.RemoteURL = &remoteURL
	repo.HTTPClient.Connection = &repository.HTTPClientConnection{
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
						resource.TestCheckResourceAttr(resName, "apt.0.distribution", repo.Apt.Distribution),
						resource.TestCheckResourceAttr(resName, "apt.0.flat", strconv.FormatBool(repo.Apt.Flat)),
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
