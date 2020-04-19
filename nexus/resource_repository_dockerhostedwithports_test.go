package nexus

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRepositoryDockerHostedWithPorts(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-docker-hosted-%s", acctest.RandString(10))
	repoOnline := true
	repoHTTPPort := acctest.RandIntRange(32767, 49152)
	repoHTTPSPort := acctest.RandIntRange(49153, 65535)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourceDockerHostedWithPorts(repoName, repoOnline, repoHTTPPort, repoHTTPSPort),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nexus_repository.docker_hosted", "name", repoName),
					resource.TestCheckResourceAttr("nexus_repository.docker_hosted", "format", "docker"),
					resource.TestCheckResourceAttr("nexus_repository.docker_hosted", "type", "hosted"),
				),
			},
			{
				ResourceName:      "nexus_repository.docker_hosted",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
				// TODO: add check for storage
				// TODO: add check for repository connectors
				// TODO: add check for Group members
				// TODO: add check for api version support
				// TODO: add tests for readonly repository
				// TODO: add tests for cleanup
			},
		},
	})
}

func createTfStmtForResourceDockerHostedWithPorts(name string, online bool, httpPort int, httpsPort int) string {
	return fmt.Sprintf(`
resource "nexus_repository" "docker_hosted" {
	name   = "%s"
	format = "docker"
	type   = "hosted"
	online = %s

	docker {
		http_port        = %d
		https_port       = %d
		force_basic_auth = true
		v1enabled        = true
	}

	storage {
		write_policy = "ALLOW"
	}
}`, name, strconv.FormatBool(online), httpPort, httpsPort)
}
