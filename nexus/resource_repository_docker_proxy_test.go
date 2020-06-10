package nexus

import (
	"fmt"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRepositoryDockerProxy(t *testing.T) {
	resName := "nexus_repository.docker_proxy"
	repoName := fmt.Sprintf("test-repo-docker-proxy-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourceDockerProxy(repoName, "https://index.docker.io/"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", repoName),
					resource.TestCheckResourceAttr(resName, "format", nexus.RepositoryFormatDocker),
					resource.TestCheckResourceAttr(resName, "type", nexus.RepositoryTypeProxy),
					resource.TestCheckResourceAttr(resName, "online", "true"),
				),
				// TODO: add check for storage
				// TODO: add check for repository connectors
				// TODO: add check for Group members
				// TODO: add check for api version support
				// TODO: add tests for readonly repository
				// TODO: add tests for cleanup
				// TODO: add tests for docker proxy specific parameters
			},
		},
	})
}

func createTfStmtForResourceDockerProxy(name string, indexURL string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "docker_proxy" {
	format = "%s"
	name   = "%s"
	online = true
	type   = "%s"

	docker {
		force_basic_auth = true
		v1enabled        = false
	}

	docker_proxy {
		index_type = "HUB"
		index_url  = "%s"
	}

	http_client {}

	negative_cache {
		enabled = true
		ttl     = 1440
	}

	proxy {
		remote_url  = "https://index.docker.io"
	}

	storage {
		blob_store_name = "default"
		write_policy    = "ALLOW"
	}
}`, nexus.RepositoryFormatDocker, name, nexus.RepositoryTypeProxy, indexURL)
}
