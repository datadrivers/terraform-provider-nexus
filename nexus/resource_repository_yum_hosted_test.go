package nexus

import (
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testAccResourceRepositoryYumHosted() repository.LegacyRepository {
	repo := testAccResourceRepositoryHosted(repository.RepositoryFormatYum)
	repo.Yum = &repository.Yum{
		DeployPolicy:  repository.YumDeployPolicyPermissive,
		RepodataDepth: 0,
	}
	return repo
}

func TestAccResourceRepositoryYumHosted(t *testing.T) {
	repo := testAccResourceRepositoryYumHosted()
	resName := testAccResourceRepositoryName(repo)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryConfig(repo),
				Check: resource.ComposeTestCheckFunc(
					resourceRepositoryTestCheckFunc(repo),
					resourceRepositoryTypeHostedTestCheckFunc(repo),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "apt.#", "0"),
						resource.TestCheckResourceAttr(resName, "bower.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker_proxy.#", "0"),
						resource.TestCheckResourceAttr(resName, "maven.#", "0"),
					),
					// Fields related to this format and type
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "yum.#", "1"),
						resource.TestCheckResourceAttr(resName, "yum.0.deploy_policy", repo.DeployPolicy),
						resource.TestCheckResourceAttr(resName, "yum.0.repodata_depth", strconv.Itoa(repo.RepodataDepth)),
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
