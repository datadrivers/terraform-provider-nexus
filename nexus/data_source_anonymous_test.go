package nexus

import (
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAnonymous(t *testing.T) {
	dataSourceName := "data.nexus_anonymous.acceptance"

	anonym := security.AnonymousAccessSettings{
		Enabled:   true,
		UserID:    acctest.RandString(20),
		RealmName: "NexusAuthenticatingRealm",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAnonymousConfig(anonym),
				Check:  nil,
			},
			{
				Config: testAccResourceAnonymousConfig(anonym) + testAccDataSourceAnonymousConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "enabled", strconv.FormatBool(anonym.Enabled)),
					resource.TestCheckResourceAttr(dataSourceName, "user_id", anonym.UserID),
					resource.TestCheckResourceAttr(dataSourceName, "realm_name", anonym.RealmName),
				),
			},
		},
	})
}

func testAccDataSourceAnonymousConfig() string {
	return `
data "nexus_anonymous" "acceptance" {}
`
}
