package other_test

import (
	"testing"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceHTTPClient(t *testing.T) {
	resName := "data.nexus_http_client.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceHTTPClientConfig("ds-proxy.example.com", 8888, "example.internal"),
			},
			{
				Config: testAccResourceHTTPClientConfig("ds-proxy.example.com", 8888, "example.internal") + testAccDataSourceHTTPClientConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "timeout", "30"),
					resource.TestCheckResourceAttr(resName, "retries", "3"),
					resource.TestCheckResourceAttr(resName, "http_proxy.0.host", "ds-proxy.example.com"),
					resource.TestCheckResourceAttr(resName, "http_proxy.0.port", "8888"),
				),
			},
		},
	})
}

func testAccDataSourceHTTPClientConfig() string {
	return `
data "nexus_http_client" "acceptance" {}
`
}
