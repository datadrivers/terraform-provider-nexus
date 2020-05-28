package nexus

import (
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceRepositoryGroup(t *testing.T) {
	repoName := "maven-public"
	resourceName := "data.nexus_repository.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRepository(repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", repoName),
						resource.TestCheckResourceAttr(resourceName, "name", repoName),
						resource.TestCheckResourceAttr(resourceName, "format", nexus.RepositoryFormatMaven2),
						resource.TestCheckResourceAttr(resourceName, "type", nexus.RepositoryTypeGroup),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "online", "true"),
						// Storage
						resource.TestCheckResourceAttr(resourceName, "storage.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "storage.0.blob_store_name", "default"),
						resource.TestCheckResourceAttr(resourceName, "storage.0.strict_content_type_validation", "true"),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "group.#", "1"),
						//resource.TestCheckResourceAttr(resourceName, "group.0", "maven-releases"),
					),
				),
			},
		},
	})
}
