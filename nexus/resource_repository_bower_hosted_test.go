package nexus

import (
	"strconv"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testAccRepositoryBowerHosted() nexus.Repository {
	repo := testAccResourceRepositoryHosted(nexus.RepositoryFormatBower)
	repo.RepositoryBower = &nexus.RepositoryBower{
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
						resource.TestCheckResourceAttr(resName, "bower.0.rewrite_package_urls", strconv.FormatBool(repo.RepositoryBower.RewritePackageUrls)),
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
