package nexus

import (
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testAccResourceRepositoryMavenProxy() nexus.Repository {
	repo := testAccResourceRepositoryProxy(nexus.RepositoryFormatMaven2)
	repo.RepositoryMaven = &nexus.RepositoryMaven{
		LayoutPolicy:  "STRICT",
		VersionPolicy: "SNAPSHOT",
	}
	repo.RepositoryProxy.RemoteURL = "https://www.example.com"
	return repo
}

func TestAccResourceRepositoryMavenProxy(t *testing.T) {
	repo := testAccResourceRepositoryMavenProxy()
	resName := testAccResourceRepositoryName(repo)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryConfig(repo),
				Check: resource.ComposeTestCheckFunc(
					resourceRepositoryTestCheckFunc(repo),
					resourceRepositoryTypeProxyTestCheckFunc(repo),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
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
				ImportStateVerifyIgnore: []string{"http_client.0.authentication.0.password"},
			},
		},
	})
}
