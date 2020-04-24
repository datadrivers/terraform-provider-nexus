package nexus

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccRepositoryDockerGroup(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-docker-group-%s", acctest.RandString(10))
	dockerProxyRepoName := fmt.Sprintf("test-repo-docker-group-member-%s", acctest.RandString(10))
	dockerHostedWithPortsRepoName := fmt.Sprintf("test-repo-docker-group-member-%s", acctest.RandString(10))
	dockerHostedWithoutPortsRepoName := fmt.Sprintf("test-repo-docker-group-member-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourceDockerHostedWithPorts(dockerHostedWithPortsRepoName) + createTfStmtForResourceDockerHostedWithoutPorts(dockerHostedWithoutPortsRepoName) + createTfStmtForResourceDockerProxy(dockerProxyRepoName) + createTfStmtForResourceDockerGroup(repoName),
				Check: resource.ComposeTestCheckFunc(
					// Base and common repo props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "id", repoName),
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "name", repoName),
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "format", "docker"),
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "type", "group"),
					),
					// Common fields
					// Online
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "online", "true"),
						// Storage
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "storage.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "storage.0.blob_store_name", "default"),
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "storage.0.strict_content_type_validation", "true"),
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "storage.0.write_policy", ""),
					),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "apt.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "apt_signing.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "bower.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "maven.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "docker_proxy.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "http_client.#", "0"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "proxy.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "negative_cache.#", "0"),
					),

					// Fields related to this format and type
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "docker.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "docker.0.force_basic_auth", "true"),
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "docker.0.http_port", "8085"),
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "docker.0.https_port", "8086"),
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "docker.0.v1enabled", "false"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "group.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.docker_group", "group.0.member_names.#", "3"),
					),
					// FIXME: (BUG) Incorrect member_names state representation.
					// For some reasons, 1st element in array is not stored as group.0.member_names.0, but instead it's stored
					// as group.0.member_names.2941663215 where 2941663215 is a "random" number.
					// This number changes from test run to test run.
					// It may be a pointer to int instead of int itself, but it's not clear and requires additional research.
					// resource.TestCheckResourceAttr("nexus_repository.docker_group", "group.0.member_names.2941663215", memberRepoName),
					// TODO: add check for repository connectors
					// TODO: add tests for readonly repository
				),
			},
		},
	})
}

func createTfStmtForResourceDockerGroup(name string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "docker_group" {
	name   = "%s"
	format = "docker"
	type   = "group"
	online = true

	group {
		member_names = [nexus_repository.docker_proxy.name, nexus_repository.docker_hosted_with_ports.name, nexus_repository.docker_hosted_without_ports.name ]
	}

	docker {
		force_basic_auth = true
		http_port        = 8085
		https_port       = 8086
		v1enabled        = false
	}

	storage {
		blob_store_name                = "default"
		strict_content_type_validation = true
	}
}`, name)
}
