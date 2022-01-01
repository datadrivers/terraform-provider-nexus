package deprecated_test

import (
	"fmt"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccResourceRepositoryNpmGroup() repository.LegacyRepository {
	repo := testAccResourceRepositoryGroup(repository.RepositoryFormatNPM)
	return repo
}

func TestAccResourceRepositoryNpmGroup(t *testing.T) {
	proxyRepo := testAccResourceRepositoryNpmProxy()
	proxyRepoResName := testAccResourceRepositoryName(proxyRepo)

	hostedRepo := testAccResourceRepositoryNpmHosted()
	hostedRepoResName := testAccResourceRepositoryName(hostedRepo)

	repo := testAccResourceRepositoryNpmGroup()
	repo.Group.MemberNames = []string{
		fmt.Sprintf("%s.name", proxyRepoResName),
		fmt.Sprintf("%s.name", hostedRepoResName),
	}
	resName := testAccResourceRepositoryName(repo)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryConfig(proxyRepo) + testAccResourceRepositoryConfig(hostedRepo) + testAccResourceRepositoryConfig(repo),
				Check: resource.ComposeTestCheckFunc(
					resourceRepositoryTestCheckFunc(repo),
					resourceRepositoryTypeGroupTestCheckFunc(repo),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "apt.#", "0"),
						resource.TestCheckResourceAttr(resName, "apt_signing.#", "0"),
						resource.TestCheckResourceAttr(resName, "bower.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker_proxy.#", "0"),
						resource.TestCheckResourceAttr(resName, "maven.#", "0"),
					),
				),
			},
			{
				ResourceName:      resName,
				ImportStateId:     repo.Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
