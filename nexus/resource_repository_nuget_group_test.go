package nexus

import (
	"fmt"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccResourceRepositoryNugetGroup() repository.LegacyRepository {
	repo := testAccResourceRepositoryGroup(repository.RepositoryFormatNuget)
	return repo
}

func TestAccResourceRepositoryNugetGroup(t *testing.T) {
	hostedRepo := testAccResourceRepositoryNugetHosted()
	hostedRepoResName := testAccResourceRepositoryName(hostedRepo)

	proxyRepo := testAccResourceRepositoryNugetProxy()
	proxyRepoResName := testAccResourceRepositoryName(proxyRepo)

	repo := testAccResourceRepositoryNugetGroup()
	repo.Group.MemberNames = []string{
		fmt.Sprintf("%s.name", hostedRepoResName),
		fmt.Sprintf("%s.name", proxyRepoResName),
	}
	resName := testAccResourceRepositoryName(repo)

	resource.Test(t, resource.TestCase{

		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryConfig(hostedRepo) + testAccResourceRepositoryConfig(proxyRepo) + testAccResourceRepositoryConfig(repo),
				Check: resource.ComposeTestCheckFunc(
					resourceRepositoryTestCheckFunc(repo),
					resourceRepositoryTypeGroupTestCheckFunc(repo),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "maven.#", "0"),
						resource.TestCheckResourceAttr(resName, "apt.#", "0"),
						resource.TestCheckResourceAttr(resName, "apt_signing.#", "0"),
						resource.TestCheckResourceAttr(resName, "bower.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker_proxy.#", "0"),
					),
				),
			},
			{
				ResourceName:            resName,
				ImportStateId:           repo.Name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"group.0.member_names."},
			},
		},
	})
}
