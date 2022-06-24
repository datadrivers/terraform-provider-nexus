package repository_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccDataSourceRepositoryP2ProxyConfig() string {
	return `
data "nexus_repository_p2_proxy" "acceptance" {
	name   = nexus_repository_p2_proxy.acceptance.id
}`
}

func TestAccDataSourceRepositoryP2Proxy(t *testing.T) {
	repoUsingDefaults := repository.P2ProxyRepository{
		Name:   fmt.Sprintf("acceptance-%s", acctest.RandString(10)),
		Online: true,
		Proxy: repository.Proxy{
			RemoteURL: "http://download.eclipse.org/releases/2019-09",
		},
		Storage: repository.Storage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: true,
		},
	}

	dataSourceName := "data.nexus_repository_p2_proxy.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRepositoryP2ProxyConfig(repoUsingDefaults) + testAccDataSourceRepositoryP2ProxyConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "id", repoUsingDefaults.Name),
						resource.TestCheckResourceAttr(dataSourceName, "name", repoUsingDefaults.Name),
						resource.TestCheckResourceAttr(dataSourceName, "online", strconv.FormatBool(repoUsingDefaults.Online)),
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
