package nexus

import (
	"fmt"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRepositoryMavenGroup(t *testing.T) {
	resName := "nexus_repository.maven_group"
	repoName := fmt.Sprintf("test-repo-maven-group-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourceMavenGroup(repoName),
				Check: resource.ComposeTestCheckFunc(
					// Base and common repo props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "id", repoName),
						resource.TestCheckResourceAttr(resName, "name", repoName),
						resource.TestCheckResourceAttr(resName, "format", nexus.RepositoryFormatMaven2),
						resource.TestCheckResourceAttr(resName, "type", nexus.RepositoryTypeGroup),
					),
					// Common fields
					// Online
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "online", "true"),
						// Storage
						resource.TestCheckResourceAttr(resName, "storage.#", "1"),
						resource.TestCheckResourceAttr(resName, "storage.0.blob_store_name", "default"),
						resource.TestCheckResourceAttr(resName, "storage.0.strict_content_type_validation", "true"),
						resource.TestCheckResourceAttr(resName, "storage.0.write_policy", ""),
					),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "apt.#", "0"),
						resource.TestCheckResourceAttr(resName, "apt_signing.#", "0"),
						resource.TestCheckResourceAttr(resName, "bower.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker_proxy.#", "0"),
						resource.TestCheckResourceAttr(resName, "http_client.#", "0"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "proxy.#", "0"),
						resource.TestCheckResourceAttr(resName, "negative_cache.#", "0"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "group.#", "1"),
						resource.TestCheckResourceAttr(resName, "group.0.member_names.#", "2"),
					),
					// FIXME: (BUG) Incorrect member_names state representation.
					// For some reasons, 1st element in array is not stored as group.0.member_names.0, but instead it's stored
					// as group.0.member_names.2941663215 where 2941663215 is a "random" number.
					// This number changes from test run to test run.
					// It may be a pointer to int instead of int itself, but it's not clear and requires additional research.
					// resource.TestCheckResourceAttr(resName, "group.0.member_names.2941663215", memberRepoName),
					// TODO: add check for repository connectors
					// TODO: add tests for readonly repository
				),
			},
		},
	})
}

func createTfStmtForResourceMavenGroup(name string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "maven_group" {
	format = "%s"
	name   = "%s"
	online = true
	type   = "%s"

	group {
		member_names = ["maven-releases", "maven-public"]
	}

	storage {
		blob_store_name                = "default"
		strict_content_type_validation = true
	}
}`, nexus.RepositoryFormatMaven2, name, nexus.RepositoryTypeGroup)
}
