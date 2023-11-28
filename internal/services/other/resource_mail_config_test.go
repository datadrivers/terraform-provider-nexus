package other_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/dre2004/go-nexus-client/nexus3/schema"
	"github.com/dre2004/terraform-provider-nexus/internal/acceptance"
	"github.com/dre2004/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceMailConfig(t *testing.T) {
	resName := "nexus_mail_config.acceptance"

	mailcfg := schema.MailConfig{
		Host:        acctest.RandString(20),
		Port:        acctest.RandIntRange(1, 65535),
		FromAddress: acctest.RandString(10) + "@" + acctest.RandString(10) + "." + acctest.RandString(3),

		Username:                      tools.GetStringPointer(acctest.RandString(10)),
		SubjectPrefix:                 tools.GetStringPointer(acctest.RandString(10)),
		Enabled:                       tools.GetBoolPointer(true),
		StartTlsEnabled:               tools.GetBoolPointer(true),
		StartTlsRequired:              tools.GetBoolPointer(true),
		SslOnConnectEnabled:           tools.GetBoolPointer(true),
		SslServerIdentityCheckEnabled: tools.GetBoolPointer(true),
		NexusTrustStoreEnabled:        tools.GetBoolPointer(true),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMailConfigConfig(mailcfg),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "host", mailcfg.Host),
					resource.TestCheckResourceAttr(resName, "from_address", mailcfg.FromAddress),
					resource.TestCheckResourceAttr(resName, "port", strconv.Itoa(mailcfg.Port)),

					resource.TestCheckResourceAttr(resName, "username", *mailcfg.Username),
					resource.TestCheckResourceAttr(resName, "subject_prefix", *mailcfg.SubjectPrefix),
					resource.TestCheckResourceAttr(resName, "enabled", strconv.FormatBool(*mailcfg.Enabled)),
					resource.TestCheckResourceAttr(resName, "start_tls_enabled", strconv.FormatBool(*mailcfg.StartTlsEnabled)),
					resource.TestCheckResourceAttr(resName, "start_tls_required", strconv.FormatBool(*mailcfg.StartTlsRequired)),
					resource.TestCheckResourceAttr(resName, "ssl_on_connect_enabled", strconv.FormatBool(*mailcfg.SslOnConnectEnabled)),
					resource.TestCheckResourceAttr(resName, "ssl_server_identity_check_enabled", strconv.FormatBool(*mailcfg.SslServerIdentityCheckEnabled)),
					resource.TestCheckResourceAttr(resName, "nexus_trust_store_enabled", strconv.FormatBool(*mailcfg.NexusTrustStoreEnabled)),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateId:     "cfg",
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceMailConfigConfig(mailcfg schema.MailConfig) string {
	return fmt.Sprintf(`
	resource "nexus_mail_config" "acceptance" {
		port                				= %d
		host                				= "%s"
		from_address        				= "%s"

		username							= "%s"
		subject_prefix      				= "%s"
		enabled								= "%v"
		start_tls_enabled   				= "%v"
		start_tls_required  				= "%v"
		ssl_on_connect_enabled   			= "%v"
		ssl_server_identity_check_enabled   = "%v"
		nexus_trust_store_enabled   		= "%v"
	  }
`, mailcfg.Port, mailcfg.Host, mailcfg.FromAddress, *mailcfg.Username, *mailcfg.SubjectPrefix, *mailcfg.Enabled, *mailcfg.StartTlsEnabled, *mailcfg.StartTlsRequired, *mailcfg.SslOnConnectEnabled, *mailcfg.SslServerIdentityCheckEnabled, *mailcfg.NexusTrustStoreEnabled)
}
