package nexus

import (
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testAccResourceRepositoryDockerProxy() nexus.Repository {
	repo := testAccResourceRepositoryProxy(nexus.RepositoryFormatDocker)
	repo.RepositoryDocker = &nexus.RepositoryDocker{}

	indexURL := "https://index.docker.io/"
	repo.RepositoryDockerProxy = &nexus.RepositoryDockerProxy{
		IndexType: "HUB",
		IndexURL:  &indexURL,
	}
	repo.RepositoryProxy.RemoteURL = "https://registry-1.docker.io"
	return repo
}

func TestAccResourceRepositoryDockerProxy(t *testing.T) {
	repo := testAccResourceRepositoryDockerProxy()
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
