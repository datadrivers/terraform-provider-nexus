package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRepositoryDockerProxy(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-docker-proxy-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourceDockerProxy(repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nexus_repository.docker_proxy", "name", repoName),
					resource.TestCheckResourceAttr("nexus_repository.docker_proxy", "format", "docker"),
					resource.TestCheckResourceAttr("nexus_repository.docker_proxy", "type", "proxy"),
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

func createTfStmtForResourceDockerProxy(name string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "docker_proxy" {
	name   = "%s"
	type   = "proxy"
	format = "docker"

	docker {
		force_basic_auth = true
		v1enabled        = false
	}

	docker_proxy {
		index_type = "HUB"
		index_url  = "http://www.example.com"
	}

	http_client {
		authentication {
			type = "username"
		}
	}

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
}`, name)
}
