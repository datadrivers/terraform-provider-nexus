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

func testAccDataSourceRepositoryDockerProxyConfig() string {
	return `
data "nexus_repository_docker_proxy" "acceptance" {
	name   = nexus_repository_docker_proxy.acceptance.id
}`
}

func TestAccDataSourceRepositoryDockerProxy(t *testing.T) {
	repoUsingDefaults := repository.DockerProxyRepository{
		Name:   fmt.Sprintf("acceptance-%s", acctest.RandString(10)),
		Online: true,
		Docker: repository.Docker{
			ForceBasicAuth: true,
			V1Enabled:      true,
		},
		DockerProxy: repository.DockerProxy{
			IndexType: repository.DockerProxyIndexTypeHub,
		},
		Proxy: repository.Proxy{
			RemoteURL: "https://registry-1.docker.io",
		},
		Storage: repository.Storage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
		},
	}

	dataSourceName := "data.nexus_repository_docker_proxy.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryDockerProxyConfig(repoUsingDefaults) + testAccDataSourceRepositoryDockerProxyConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "id", repoUsingDefaults.Name),
						resource.TestCheckResourceAttr(dataSourceName, "name", repoUsingDefaults.Name),
						resource.TestCheckResourceAttr(dataSourceName, "online", strconv.FormatBool(repoUsingDefaults.Online)),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "docker.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "docker.0.force_basic_auth", strconv.FormatBool(repoUsingDefaults.Docker.ForceBasicAuth)),
						resource.TestCheckResourceAttr(dataSourceName, "docker.0.v1_enabled", strconv.FormatBool(repoUsingDefaults.Docker.V1Enabled)),
						resource.TestCheckResourceAttr(dataSourceName, "docker_proxy.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "docker_proxy.0.index_type", string(repoUsingDefaults.DockerProxy.IndexType)),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "http_client.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "http_client.0.authentication.#", "0"),
						resource.TestCheckResourceAttr(dataSourceName, "http_client.0.connection.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "negative_cache.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "proxy.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "proxy.0.remote_url", repoUsingDefaults.Proxy.RemoteURL),
						resource.TestCheckResourceAttr(dataSourceName, "storage.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "storage.0.blob_store_name", repoUsingDefaults.Storage.BlobStoreName),
						resource.TestCheckResourceAttr(dataSourceName, "storage.0.strict_content_type_validation", strconv.FormatBool(repoUsingDefaults.Storage.StrictContentTypeValidation)),
					),
				),
			},
		},
	})
}

func TestAccProDataSourceRepositoryDockerProxy(t *testing.T) {
	if tools.GetEnv("SKIP_PRO_TESTS", "false") == "true" {
		t.Skip("Skipping Nexus Pro Tests")
	}
	name := fmt.Sprintf("acceptance-%s", acctest.RandString(10))
	repoUsingDefaults := repository.DockerProxyRepository{
		Name:   name,
		Online: true,
		Docker: repository.Docker{
			ForceBasicAuth: true,
			V1Enabled:      true,
			Subdomain:      tools.GetStringPointer(name),
		},
		DockerProxy: repository.DockerProxy{
			IndexType: repository.DockerProxyIndexTypeHub,
		},
		Proxy: repository.Proxy{
			RemoteURL: "https://registry-1.docker.io",
		},
		Storage: repository.Storage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
		},
	}

	dataSourceName := "data.nexus_repository_docker_proxy.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryDockerProxyConfig(repoUsingDefaults) + testAccDataSourceRepositoryDockerProxyConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "id", repoUsingDefaults.Name),
						resource.TestCheckResourceAttr(dataSourceName, "name", repoUsingDefaults.Name),
						resource.TestCheckResourceAttr(dataSourceName, "online", strconv.FormatBool(repoUsingDefaults.Online)),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "docker.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "docker.0.force_basic_auth", strconv.FormatBool(repoUsingDefaults.Docker.ForceBasicAuth)),
						resource.TestCheckResourceAttr(dataSourceName, "docker.0.v1_enabled", strconv.FormatBool(repoUsingDefaults.Docker.V1Enabled)),
						resource.TestCheckResourceAttr(dataSourceName, "docker.0.subdomain", string(*repo.Docker.Subdomain)),
						resource.TestCheckResourceAttr(dataSourceName, "docker_proxy.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "docker_proxy.0.index_type", string(repoUsingDefaults.DockerProxy.IndexType)),
					),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "http_client.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "http_client.0.authentication.#", "0"),
						resource.TestCheckResourceAttr(dataSourceName, "http_client.0.connection.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "negative_cache.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "proxy.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "proxy.0.remote_url", repoUsingDefaults.Proxy.RemoteURL),
						resource.TestCheckResourceAttr(dataSourceName, "storage.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "storage.0.blob_store_name", repoUsingDefaults.Storage.BlobStoreName),
						resource.TestCheckResourceAttr(dataSourceName, "storage.0.strict_content_type_validation", strconv.FormatBool(repoUsingDefaults.Storage.StrictContentTypeValidation)),
					),
				),
			},
		},
	})
}
