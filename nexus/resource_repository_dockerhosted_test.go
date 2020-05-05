package nexus

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccRepositoryDockerHostedWithPorts(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-docker-hosted-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourceDockerHostedWithPorts(repoName),
				Check: resource.ComposeTestCheckFunc(
					// Base and common repo props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "id", repoName),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "name", repoName),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "format", "docker"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "type", "hosted"),
					),
					// Common fields
					// Online
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "online", "true"),
						// Storage
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "storage.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "storage.0.blob_store_name", "default"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "storage.0.strict_content_type_validation", "true"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "storage.0.write_policy", "ALLOW"),
					),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "apt.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "apt_signing.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "bower.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "maven.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "docker_proxy.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "http_client.#", "0"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "proxy.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "negative_cache.#", "0"),
					),

					// Fields related to this format and type
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "docker.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "docker.0.force_basic_auth", "true"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "docker.0.v1enabled", "true"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "docker.0.http_port", "8083"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_with_ports", "docker.0.https_port", "8084"),
					),
				),
			},
			{
				ResourceName:      "nexus_repository.docker_hosted_with_ports",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func createTfStmtForResourceDockerHostedWithPorts(name string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "docker_hosted_with_ports" {
	name   = "%s"
	format = "docker"
	type   = "hosted"
	online = true

	docker {
		http_port        = 8083
		https_port       = 8084
		force_basic_auth = true
		v1enabled        = true
	}

	storage {
		write_policy = "ALLOW"
	}
}`, name)
}

func TestAccRepositoryDockerHostedWithoutPorts(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-docker-hosted-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourceDockerHostedWithoutPorts(repoName),
				Check: resource.ComposeTestCheckFunc(
					// Base and common repo props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "id", repoName),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "name", repoName),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "format", "docker"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "type", "hosted"),
					),
					// Common fields
					// Online
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "online", "true"),
						// Storage
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "storage.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "storage.0.blob_store_name", "default"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "storage.0.strict_content_type_validation", "true"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "storage.0.write_policy", "ALLOW"),
					),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "apt.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "apt_signing.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "bower.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "maven.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "docker_proxy.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "http_client.#", "0"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "proxy.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "negative_cache.#", "0"),
					),

					// Fields related to this format and type
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "docker.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "docker.0.force_basic_auth", "true"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "docker.0.v1enabled", "true"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "docker.0.http_port", "0"),
						resource.TestCheckResourceAttr("nexus_repository.docker_hosted_without_ports", "docker.0.https_port", "0"),
					),
				),
			},
			{
				ResourceName:      "nexus_repository.docker_hosted_without_ports",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func createTfStmtForResourceDockerHostedWithoutPorts(name string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "docker_hosted_without_ports" {
	name   = "%s"
	format = "docker"
	type   = "hosted"
	online = true

	docker {
		force_basic_auth = true
		v1enabled        = true
	}

	storage {
		write_policy = "ALLOW"
	}
}`, name)
}
