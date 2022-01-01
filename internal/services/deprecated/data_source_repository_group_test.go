package deprecated_test

import (
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRepositoryGroup(t *testing.T) {
	repoName := "maven-public"
	resourceName := "data.nexus_repository.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRepositoryConfig(repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", repoName),
						resource.TestCheckResourceAttr(resourceName, "name", repoName),
						resource.TestCheckResourceAttr(resourceName, "format", repository.RepositoryFormatMaven2),
						resource.TestCheckResourceAttr(resourceName, "type", repository.RepositoryTypeGroup),
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
