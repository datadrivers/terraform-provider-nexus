package nexus

import (
	"fmt"
	"strconv"
	"testing"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccResourceSecurityLDAP(t *testing.T) {
	resName := "nexus_security_ldap.acceptance"
	ldap := nexus.LDAP{
		AuthPassword:                "1234567890",
		AuthSchema:                  "SIMPLE",
		AuthUserName:                "admin",
		ConnectionRetryDelaySeconds: uint(1),
		ConnectionTimeoutSeconds:    uint(1),
		GroupType:                   "static",
		Host:                        "127.0.0.1",
		MaxIncidentCount:            uint(1),
		Name:                        "ci-test",
		Port:                        389,
		Protocol:                    "LDAP",
		SearchBase:                  "dc=example,dc=com",
		UserEmailAddressAttribute:   "mail",
		UserIDAttribute:             "uid",
		UserObjectClass:             "inetOrgPerson",
		UserRealNameAttribute:       "cn",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityLDAP(ldap),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "auth_password"),
					resource.TestCheckResourceAttr(resName, "auth_realm", ldap.AuthRealm),
					resource.TestCheckResourceAttr(resName, "auth_schema", ldap.AuthSchema),
					resource.TestCheckResourceAttr(resName, "auth_username", ldap.AuthUserName),
					resource.TestCheckResourceAttr(resName, "connection_retry_delay_seconds", strconv.Itoa(int(ldap.ConnectionRetryDelaySeconds))),
					resource.TestCheckResourceAttr(resName, "connection_timeout_seconds", strconv.Itoa(int(ldap.ConnectionTimeoutSeconds))),
					resource.TestCheckResourceAttr(resName, "group_base_dn", ldap.GroupBaseDn),
					resource.TestCheckResourceAttr(resName, "group_id_attribute", ldap.GroupIDAttribute),
					resource.TestCheckResourceAttr(resName, "group_member_attribute", ldap.GroupMemberAttribute),
					resource.TestCheckResourceAttr(resName, "group_member_format", ldap.GroupMemberFormat),
					resource.TestCheckResourceAttr(resName, "group_object_class", ldap.GroupObjectClass),
					resource.TestCheckResourceAttr(resName, "group_subtree", strconv.FormatBool(ldap.GroupSubtree)),
					resource.TestCheckResourceAttr(resName, "host", ldap.Host),
					resource.TestCheckResourceAttr(resName, "ldap_groups_as_roles", strconv.FormatBool(ldap.LDAPGroupsAsRoles)),
					resource.TestCheckResourceAttr(resName, "max_incident_count", strconv.Itoa(int(ldap.MaxIncidentCount))),
					resource.TestCheckResourceAttr(resName, "name", ldap.Name),
					resource.TestCheckResourceAttr(resName, "port", strconv.Itoa(int(ldap.Port))),
					resource.TestCheckResourceAttr(resName, "protocol", ldap.Protocol),
					resource.TestCheckResourceAttr(resName, "search_base", ldap.SearchBase),
					resource.TestCheckResourceAttr(resName, "user_base_dn", ldap.UserBaseDN),
					resource.TestCheckResourceAttr(resName, "user_email_address_attribute", ldap.UserEmailAddressAttribute),
					resource.TestCheckResourceAttr(resName, "user_id_attribute", ldap.UserIDAttribute),
					resource.TestCheckResourceAttr(resName, "user_ldap_filter", ldap.UserLDAPFilter),
					resource.TestCheckResourceAttr(resName, "user_member_of_attribute", ldap.UserMemberOfAttribute),
					resource.TestCheckResourceAttr(resName, "user_object_class", ldap.UserObjectClass),
					resource.TestCheckResourceAttr(resName, "user_password_attribute", ldap.UserPasswordAttribute),
					resource.TestCheckResourceAttr(resName, "user_real_name_attribute", ldap.UserRealNameAttribute),
					resource.TestCheckResourceAttr(resName, "user_subtree", strconv.FormatBool(ldap.UserSubtree)),
				),
			},
			{
				ResourceName:      resName,
				ImportStateId:     ldap.Name,
				ImportState:       true,
				ImportStateVerify: true,
				// auth_password and group_type are not returned
				ImportStateVerifyIgnore: []string{"auth_password", "group_type"},
			},
		},
	})
}

func testAccResourceSecurityLDAP(ldap nexus.LDAP) string {
	return fmt.Sprintf(`
resource "nexus_security_ldap" "acceptance" {
	auth_password                  = "%s"
	auth_schema                    = "%s"
	auth_username                  = "%s"
	connection_retry_delay_seconds = %d
	connection_timeout_seconds     = %d
	group_type                     = "%s"
	host                           = "%s"
	max_incident_count             = %d
	name                           = "%s"
	port                           = %d
	protocol                       = "%s"
	search_base                    = "%s"
	user_email_address_attribute   = "%s"
	user_id_attribute              = "%s"
	user_object_class              = "%s"
	user_real_name_attribute       = "%s"
}
`, ldap.AuthPassword, ldap.AuthSchema, ldap.AuthUserName, ldap.ConnectionRetryDelaySeconds, ldap.ConnectionTimeoutSeconds, ldap.GroupType, ldap.Host, ldap.MaxIncidentCount, ldap.Name, ldap.Port, ldap.Protocol, ldap.SearchBase, ldap.UserEmailAddressAttribute, ldap.UserIDAttribute, ldap.UserObjectClass, ldap.UserRealNameAttribute)
}
