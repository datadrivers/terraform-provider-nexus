package nexus

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRepositoryDockerHostedWithoutPorts(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-docker-hosted-%s", acctest.RandString(10))
	repoOnline := true

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRepositoryDockerHostedWithoutPorts(repoName, repoOnline),
			},
		},
	})
}
func testAccRepositoryDockerHostedWithoutPorts(name string, online bool) string {
	return fmt.Sprintf(`
resource "nexus_repository" "docker_hosted" {
	name   = "%s"
	format = "docker"
	type   = "hosted"
	online = %s

	docker {
		force_basic_auth = true
		v1enabled        = true
	}

	storage {
		write_policy = "ALLOW"
	}
}`, name, strconv.FormatBool(online))
}
