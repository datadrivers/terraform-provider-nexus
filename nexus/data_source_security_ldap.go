/*
Use this data source to work with LDAP security

Example Usage

```hcl
data "nexus_security_ldap" "default" {}
```
*/
package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSecurityLDAP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecurityLDAPRead,

		Schema: map[string]*schema.Schema{
			"ldap": {
				Computed:    true,
				Description: "List of ldap configrations",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth_password": {
							Computed:    true,
							Description: "The password to bind with",
							Type:        schema.TypeString,
						},
						"auth_realm": {
							Computed:    true,
							Description: "The SASL realm to bind to",
							Type:        schema.TypeString,
						},
						"auth_schema": {
							Computed:    true,
							Description: "Authentication scheme used for connecting to LDAP server",
							Type:        schema.TypeString,
						},
						"auth_username": {
							Computed:    true,
							Description: "This must be a fully qualified username if simple authentication is used",
							Type:        schema.TypeString,
						},
						"id": {
							Computed:    true,
							Description: "The if of the ldap configuration",
							Type:        schema.TypeString,
						},
						"connection_retry_delay_seconds": {
							Computed:    true,
							Description: "How long to wait before retrying",
							Type:        schema.TypeInt,
						},
						"connection_timeout_seconds": {
							Computed:    true,
							Description: "How long to wait before timeout",
							Type:        schema.TypeInt,
						},
						"group_base_dn": {
							Computed:    true,
							Description: "The relative DN where group objects are found (e.g. ou=Group)",
							Type:        schema.TypeString,
						},
						"group_id_attribute": {
							Computed:    true,
							Description: "This field specifies the attribute of the Object class that defines the Group ID",
							Type:        schema.TypeString,
						},
						"group_member_attribute": {
							Computed:    true,
							Description: "LDAP attribute containing the usernames for the group",
							Type:        schema.TypeString,
						},
						"group_member_format": {
							Computed:    true,
							Description: "The format of user ID stored in the group member attribute",
							Type:        schema.TypeString,
						},
						"group_object_class": {
							Computed:    true,
							Description: "LDAP class for group objects",
							Type:        schema.TypeString,
						},
						"group_subtree": {
							Computed:    true,
							Description: "Are groups located in structures below the group base DN",
							Type:        schema.TypeString,
						},
						"group_type": {
							Computed:    true,
							Description: "Defines a type of groups used: static (a group contains a list of users) or dynamic (a user contains a list of groups)",
							Type:        schema.TypeString,
						},
						"host": {
							Computed:    true,
							Description: "LDAP server connection hostname",
							Type:        schema.TypeString,
						},
						"ldap_groups_as_roles": {
							Computed:    true,
							Description: "Denotes whether LDAP assigned roles are used as Nexus Repository Manager roles",
							Type:        schema.TypeBool,
						},
						"max_incident_count": {
							Computed:    true,
							Description: "How many retry attempts",
							Type:        schema.TypeInt,
						},
						"name": {
							Computed:    true,
							Description: "LDAP server name",
							Type:        schema.TypeString,
						},
						"port": {
							Computed:    true,
							Description: "LDAP server connection port to use",
							Type:        schema.TypeInt,
						},
						"protocol": {
							Computed:    true,
							Description: "LDAP server connection Protocol to use",
							Type:        schema.TypeString,
						},
						"search_base": {
							Computed:    true,
							Description: "LDAP location to be added to the connection URL",
							Type:        schema.TypeString,
						},
						"use_trust_store": {
							Computed:    true,
							Description: "Whether to use certificates stored in Nexus Repository Manager's truststore",
							Type:        schema.TypeBool,
						},
						"user_base_dn": {
							Computed:    true,
							Description: "The relative DN where user objects are found (e.g. ou=people). This value will have the Search base DN value appended to form the full User search base DN.",
							Type:        schema.TypeString,
						},
						"user_email_address_attribute": {
							Computed:    true,
							Description: "This is used to find an email address given the user ID",
							Type:        schema.TypeString,
						},
						"user_id_attribute": {
							Computed:    true,
							Description: "This is used to find a user given its user ID",
							Type:        schema.TypeString,
						},
						"user_ldap_filter": {
							Computed:    true,
							Description: "LDAP search filter to limit user search",
							Type:        schema.TypeString,
						},
						"user_member_of_attribute": {
							Computed:    true,
							Description: "Set this to the attribute used to store the attribute which holds groups DN in the user object",
							Type:        schema.TypeString,
						},
						"user_object_class": {
							Computed:    true,
							Description: "LDAP class for user objects",
							Type:        schema.TypeString,
						},
						"user_password_attribute": {
							Computed:    true,
							Description: "If this field is blank the user will be authenticated against a bind with the LDAP server",
							Type:        schema.TypeString,
						},
						"user_real_name_attribute": {
							Computed:    true,
							Description: "This is used to find a real name given the user ID",
							Type:        schema.TypeString,
						},
						"user_subtree": {
							Computed:    true,
							Description: "Are users located in structures below the user base DN?",
							Type:        schema.TypeBool,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceSecurityLDAPRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	ldapServer, err := client.Security.LDAP.List()
	if err != nil {
		return err
	}

	d.SetId("ldap")
	if err := d.Set("ldap", flattenSecurityLDAP(ldapServer)); err != nil {
		return err
	}

	return nil
}

func flattenSecurityLDAP(ldap []security.LDAP) []map[string]interface{} {
	if ldap == nil {
		return nil
	}
	data := make([]map[string]interface{}, len(ldap))
	for i, server := range ldap {
		data[i] = map[string]interface{}{
			"auth_password":                  server.AuthPassword,
			"auth_realm":                     server.AuthRealm,
			"auth_schema":                    server.AuthSchema,
			"auth_username":                  server.AuthUserName,
			"connection_retry_delay_seconds": server.ConnectionRetryDelaySeconds,
			"connection_timeout_seconds":     server.ConnectionTimeoutSeconds,
			"group_base_dn":                  server.GroupBaseDn,
			"group_id_attribute":             server.GroupIDAttribute,
			"group_member_attribute":         server.GroupMemberAttribute,
			"group_member_format":            server.GroupMemberFormat,
			"group_object_class":             server.GroupObjectClass,
			"group_subtree":                  server.GroupSubtree,
			"group_type":                     server.GroupType,
			"host":                           server.Host,
			"id":                             server.ID,
			"ldap_groups_as_roles":           server.LDAPGroupsAsRoles,
			"max_incident_count":             server.MaxIncidentCount,
			"name":                           server.Name,
			"port":                           server.Port,
			"protocol":                       server.Protocol,
			"search_base":                    server.SearchBase,
			"user_base_dn":                   server.UserBaseDN,
			"user_email_address_attribute":   server.UserEmailAddressAttribute,
			"user_id_attribute":              server.UserIDAttribute,
			"user_ldap_filter":               server.UserLDAPFilter,
			"user_member_of_attribute":       server.UserMemberOfAttribute,
			"user_object_class":              server.UserObjectClass,
			"user_password_attribute":        server.UserPasswordAttribute,
			"user_real_name_attribute":       server.UserRealNameAttribute,
			"user_subtree":                   server.UserSubtree,
			"user_trust_store":               server.UseTrustStore,
		}
	}
	return nil
}
