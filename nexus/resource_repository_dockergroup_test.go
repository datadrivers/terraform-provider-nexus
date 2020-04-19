package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRepositoryDockerGroup(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-docker-group-%s", acctest.RandString(10))
	memberRepoName := fmt.Sprintf("test-repo-docker-group-member-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourceDockerProxy(memberRepoName) + createTfStmtForResourceDockerGroup(repoName, memberRepoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nexus_repository.docker_group", "name", repoName),
					resource.TestCheckResourceAttr("nexus_repository.docker_group", "format", "docker"),
					resource.TestCheckResourceAttr("nexus_repository.docker_group", "type", "group"),
				// TODO: add check for storage
				// TODO: add check for repository connectors
				// TODO: add check for Group members
				// TODO: add check for api version support
				// TODO: add tests for readonly repository
				),
			},
		},
	})
}

func createTfStmtForResourceDockerGroup(name string, memberRepoName string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "docker_group" {
	name   = "%s"
	format = "docker"
	type   = "group"
	online = true
	
	group {
		member_names = [nexus_repository.docker_proxy.name]
	}
	
	docker {
		force_basic_auth = true
		http_port        = 8082
		https_port       = 0
		v1enabled        = false
	}
	
	storage {
		blob_store_name                = "default"
		strict_content_type_validation = true
	}
}`, name)
}
