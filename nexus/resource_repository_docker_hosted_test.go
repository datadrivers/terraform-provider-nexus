package nexus

import (
	"strconv"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testAccResourceRepositoryDockerHostedWithPorts() nexus.Repository {
	repo := testAccResourceRepositoryHosted(nexus.RepositoryFormatDocker)

	httpPort := 8085
	httpsPort := 8086

	repo.RepositoryDocker = &nexus.RepositoryDocker{
		ForceBasicAuth: false,
		HTTPPort:       &httpPort,
		HTTPSPort:      &httpsPort,
		V1Enabled:      false,
	}
	return repo
}

func TestAccResourceRepositoryDockerHostedWithPorts(t *testing.T) {
	repo := testAccResourceRepositoryDockerHostedWithPorts()
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
					// No fields related to other repo types
					// Format
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
						resource.TestCheckResourceAttr(resName, "docker.0.v1enabled", strconv.FormatBool(repo.RepositoryDocker.V1Enabled)),
						resource.TestCheckResourceAttr(resName, "docker.0.http_port", strconv.Itoa(*repo.RepositoryDocker.HTTPPort)),
						resource.TestCheckResourceAttr(resName, "docker.0.https_port", strconv.Itoa(*repo.RepositoryDocker.HTTPSPort)),
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

func testAccResourceRepositoryDockerHostedWithoutPorts() nexus.Repository {
	repo := testAccResourceRepositoryHosted(nexus.RepositoryFormatDocker)
	repo.RepositoryDocker = &nexus.RepositoryDocker{
		ForceBasicAuth: true,
		V1Enabled:      false,
	}
	return repo
}

func TestAccResourceRepositoryDockerHostedWithoutPorts(t *testing.T) {
	repo := testAccResourceRepositoryDockerHostedWithoutPorts()
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
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "apt.#", "0"),
						resource.TestCheckResourceAttr(resName, "apt_signing.#", "0"),
						resource.TestCheckResourceAttr(resName, "bower.#", "0"),
						resource.TestCheckResourceAttr(resName, "maven.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker_proxy.#", "0"),
						resource.TestCheckResourceAttr(resName, "http_client.#", "0"),
					),
					// Fields related to this format and type
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "docker.#", "1"),
						resource.TestCheckResourceAttr(resName, "docker.0.force_basic_auth", strconv.FormatBool(repo.RepositoryDocker.ForceBasicAuth)),
						resource.TestCheckResourceAttr(resName, "docker.0.v1enabled", strconv.FormatBool(repo.RepositoryDocker.V1Enabled)),
						resource.TestCheckResourceAttr(resName, "docker.0.http_port", "0"),
						resource.TestCheckResourceAttr(resName, "docker.0.https_port", "0"),
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
