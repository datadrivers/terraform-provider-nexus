package other_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceHTTPClient(t *testing.T) {
	resName := "nexus_http_client.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceHTTPClientConfig("proxy.example.com", 3128, "localhost"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "timeout", "30"),
					resource.TestCheckResourceAttr(resName, "retries", "3"),
					resource.TestCheckResourceAttr(resName, "user_agent", "terraform-provider-nexus-acc"),
					resource.TestCheckResourceAttr(resName, "http_proxy.0.enabled", "true"),
					resource.TestCheckResourceAttr(resName, "http_proxy.0.host", "proxy.example.com"),
					resource.TestCheckResourceAttr(resName, "http_proxy.0.port", "3128"),
					resource.TestCheckResourceAttr(resName, "non_proxy_hosts.#", "1"),
					resource.TestCheckTypeSetElemAttr(resName, "non_proxy_hosts.*", "localhost"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateId:     "http",
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceHTTPClientConfig("proxy2.example.com", 8080, "127.0.0.1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "http_proxy.0.host", "proxy2.example.com"),
					resource.TestCheckResourceAttr(resName, "http_proxy.0.port", strconv.Itoa(8080)),
					resource.TestCheckTypeSetElemAttr(resName, "non_proxy_hosts.*", "127.0.0.1"),
				),
			},
		},
	})
}

func testAccResourceHTTPClientConfig(host string, port int, bypass string) string {
	return fmt.Sprintf(`
resource "nexus_http_client" "acceptance" {
  timeout    = 30
  retries    = 3
  user_agent = "terraform-provider-nexus-acc"

  non_proxy_hosts = ["%s"]

  http_proxy {
    enabled = true
    host    = "%s"
    port    = %d
  }
}
`, bypass, host, port)
}
