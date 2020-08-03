package nexus

import (
	"fmt"
	"strconv"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testAccResourceRepositoryDockerGroup() nexus.Repository {
	httpPort := 8087
	httpsPort := 8088

	repo := testAccResourceRepositoryGroup(nexus.RepositoryFormatDocker)
	repo.RepositoryDocker = &nexus.RepositoryDocker{
		ForceBasicAuth: true,
		HTTPPort:       &httpPort,
		HTTPSPort:      &httpsPort,
		V1Enabled:      false,
	}
	return repo
}

func TestAccResourceRepositoryDockerGroup(t *testing.T) {
	hostedRepo := testAccResourceRepositoryDockerHostedWithPorts()
	hostedRepoResName := testAccResourceRepositoryName(hostedRepo)

	proxyRepo := testAccResourceRepositoryDockerProxy()
	proxyRepoResName := testAccResourceRepositoryName(proxyRepo)

	repo := testAccResourceRepositoryDockerGroup()
	repo.RepositoryGroup.MemberNames = []string{
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
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "apt.#", "0"),
						resource.TestCheckResourceAttr(resName, "apt_signing.#", "0"),
						resource.TestCheckResourceAttr(resName, "bower.#", "0"),
						resource.TestCheckResourceAttr(resName, "maven.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker_proxy.#", "0"),
					),
					// Fields related to this format and type
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "docker.#", "1"),
						resource.TestCheckResourceAttr(resName, "docker.0.force_basic_auth", strconv.FormatBool(repo.RepositoryDocker.ForceBasicAuth)),
						resource.TestCheckResourceAttr(resName, "docker.0.http_port", strconv.Itoa(*repo.RepositoryDocker.HTTPPort)),
						resource.TestCheckResourceAttr(resName, "docker.0.https_port", strconv.Itoa(*repo.RepositoryDocker.HTTPSPort)),
						resource.TestCheckResourceAttr(resName, "docker.0.v1enabled", strconv.FormatBool(repo.RepositoryDocker.V1Enabled)),
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
