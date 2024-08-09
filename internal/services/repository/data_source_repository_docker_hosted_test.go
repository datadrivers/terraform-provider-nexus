package repository_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccDataSourceRepositoryDockerHostedConfig() string {
	return `
data "nexus_repository_docker_hosted" "acceptance" {
	name   = nexus_repository_docker_hosted.acceptance.id
}`
}

func TestAccDataSourceRepositoryDockerHosted(t *testing.T) {
	name := fmt.Sprintf("acceptance-%s", acctest.RandString(10))
	subdomain := ""
	if tools.GetEnv("SKIP_PRO_TESTS", "false") == "false" {
		subdomain = name
	}
	repo := repository.DockerHostedRepository{
		Name:   name,
		Online: true,
		Storage: repository.DockerHostedStorage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: false,
			WritePolicy:                 "ALLOW",
		},
		Docker: repository.Docker{
			ForceBasicAuth: true,
			V1Enabled:      true,
		},
	}
	repo.Docker.Subdomain = &subdomain
	dataSourceName := "data.nexus_repository_docker_hosted.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryDockerHostedConfig(repo) + testAccDataSourceRepositoryDockerHostedConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "id", repo.Name),
						resource.TestCheckResourceAttr(dataSourceName, "name", repo.Name),
						resource.TestCheckResourceAttr(dataSourceName, "online", strconv.FormatBool(repo.Online)),
						resource.TestCheckResourceAttr(dataSourceName, "docker.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "docker.0.force_basic_auth", strconv.FormatBool(repo.Docker.ForceBasicAuth)),
						resource.TestCheckResourceAttr(dataSourceName, "docker.0.v1_enabled", strconv.FormatBool(repo.Docker.V1Enabled)),
						resource.TestCheckResourceAttr(dataSourceName, "storage.0.blob_store_name", repo.Storage.BlobStoreName),
						resource.TestCheckResourceAttr(dataSourceName, "storage.0.strict_content_type_validation", strconv.FormatBool(repo.Storage.StrictContentTypeValidation)),
						resource.TestCheckResourceAttr(dataSourceName, "storage.0.write_policy", string(repo.Storage.WritePolicy)),
					),
				),
			},
		},
	})
}
