package nexus

import (
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testAccResourceDockerProxy() repository.LegacyRepository {
	repo := testAccResourceRepositoryProxy(repository.RepositoryFormatDocker)
	repo.Docker = &repository.Docker{}

	indexURL := "https://index.docker.io/"
	repo.DockerProxy = &repository.DockerProxy{
		IndexType: "HUB",
		IndexURL:  &indexURL,
	}
	repo.Proxy.RemoteURL = "https://registry-1.docker.io"
	return repo
}

func TestAccResourceDockerProxy(t *testing.T) {
	repo := testAccResourceDockerProxy()
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
