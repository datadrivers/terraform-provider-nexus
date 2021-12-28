package nexus

import (
	"fmt"
	"strconv"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testAccResourceRepositoryYumHosted() nexus.Repository {
	repo := testAccResourceRepositoryHosted(nexus.RepositoryFormatYum)
	repo.RepositoryYum = &nexus.RepositoryYum{
		DeployPolicy:  "PERMISSIVE",
		RepodataDepth: 0,
	}
	return repo
}

func resourceYumRepositoryTestCheckFunc(repo nexus.Repository) resource.TestCheckFunc {
	resName := fmt.Sprintf("nexus_repository_yum_hosted.%s", repo.Name)
	return resource.ComposeAggregateTestCheckFunc(
		resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resName, "id", repo.Name),
			resource.TestCheckResourceAttr(resName, "name", repo.Name),
			resource.TestCheckResourceAttr(resName, "online", strconv.FormatBool(repo.Online)),
		),
		resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resName, "storage.#", "1"),
			resource.TestCheckResourceAttr(resName, "storage.0.blob_store_name", repo.RepositoryStorage.BlobStoreName),
			resource.TestCheckResourceAttr(resName, "storage.0.strict_content_type_validation", strconv.FormatBool(repo.RepositoryStorage.StrictContentTypeValidation)),
		),
	)
}

func resourceYumRepositoryTypeHostedTestCheckFunc(repo nexus.Repository) resource.TestCheckFunc {
	resName := fmt.Sprintf("nexus_repository_yum_hosted.%s", repo.Name)
	return resource.ComposeAggregateTestCheckFunc(
		resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resName, "http_client.#", "0"),
			resource.TestCheckResourceAttr(resName, "group.#", "0"),
			resource.TestCheckResourceAttr(resName, "negative_cache.#", "0"),
			resource.TestCheckResourceAttr(resName, "proxy.#", "0"),
		),
		resource.TestCheckResourceAttr(resName, "storage.0.write_policy", *repo.RepositoryStorage.WritePolicy),
	)
}

func TestAccResourceRepositoryYumHosted(t *testing.T) {
	repo := testAccResourceRepositoryYumHosted()
	resName := fmt.Sprintf("nexus_repository_yum_hosted.%s", repo.Name)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryConfig(repo),
				Check: resource.ComposeTestCheckFunc(
					resourceYumRepositoryTestCheckFunc(repo),
					resourceYumRepositoryTypeHostedTestCheckFunc(repo),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resName, "apt.#", "0"),
						resource.TestCheckResourceAttr(resName, "bower.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker.#", "0"),
						resource.TestCheckResourceAttr(resName, "docker_proxy.#", "0"),
						resource.TestCheckResourceAttr(resName, "maven.#", "0"),
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
