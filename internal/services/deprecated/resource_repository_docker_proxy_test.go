package deprecated_test

import (
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccResourceDockerProxy() repository.LegacyRepository {
	repo := testAccResourceRepositoryProxy(repository.RepositoryFormatDocker)
	repo.Docker = &repository.Docker{}

	indexURL := "https://index.docker.io/"
	indexType := repository.DockerProxyIndexTypeHub
	remoteURL := "https://registry-1.docker.io"
	repo.DockerProxy = &repository.DockerProxy{
		IndexType: &indexType,
		IndexURL:  &indexURL,
	}
	repo.Proxy.RemoteURL = &remoteURL
	return repo
}

func TestAccResourceDockerProxy(t *testing.T) {
	repo := testAccResourceDockerProxy()
	resName := testAccResourceRepositoryName(repo)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
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
