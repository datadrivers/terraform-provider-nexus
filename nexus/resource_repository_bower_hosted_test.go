package nexus

import (
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testAccRepositoryBowerHosted() repository.LegacyRepository {
	repo := testAccResourceRepositoryHosted(repository.RepositoryFormatBower)
	repo.Bower = &repository.Bower{
		RewritePackageUrls: true,
	}
	return repo
}

func TestAccResourceRepositoryBowerHosted(t *testing.T) {
	repo := testAccRepositoryBowerHosted()
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
						resource.TestCheckResourceAttr(resName, "apt_signing.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker_proxy.#", "0"),
						resource.TestCheckResourceAttr(resName, "maven.#", "0"),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "bower.#", "1"),
						resource.TestCheckResourceAttr(resName, "bower.0.rewrite_package_urls", strconv.FormatBool(repo.Bower.RewritePackageUrls)),
					),
				),
			},
			{
				ResourceName:      resName,
				ImportStateId:     repo.Name,
				ImportState:       true,
				ImportStateVerify: true,
				// TODO: verify bower configuration, bower attribute is not returned by API currently
				ImportStateVerifyIgnore: []string{"bower"},
				// TODO: add tests for readonly repository
			},
		},
	})
}
