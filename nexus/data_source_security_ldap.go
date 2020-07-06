package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSecurityLDAP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecurityLDAPRead,

		Schema: map[string]*schema.Schema{
			"ldap": {
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth_password": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"auth_realm": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"auth_schema": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"auth_username": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"connection_retry_delay_seconds": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"connection_timeout_seconds": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"group_base_dn": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"group_id_attribute": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"group_member_attribute": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"group_member_format": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"group_object_class": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"group_subtree": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"group_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"host": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"ldap_groups_as_roles": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"max_incident_count": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"port": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"protocol": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"search_base": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"use_base_con": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"use_subtree": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"use_trust_store": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"user_email_address_attribute": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"user_id_attribute": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"user_ldap_filter": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"user_member_of_attribute": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"user_object_class": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"user_password_attribute": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"user_real_name_attribute": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceSecurityLDAPRead(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	ldapServer, err := nexusClient.LDAPList()
	if err != nil {
		return err
	}

	d.SetId("ldap")
	if err := d.Set("ldap", flattenSecurityLDAP(ldapServer)); err != nil {
		return err
	}

	return nil
}

func flattenSecurityLDAP(ldap []nexus.LDAP) []map[string]interface{} {
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
			"id":                             server.ID,
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
			"ldap_groups_as_roles":           server.LDAPGroupsAsRoles,
			"max_incident_count":             server.MaxIncidentCount,
			"name":                           server.Name,
			"port":                           server.Port,
			"protocol":                       server.Protocol,
			"search_base":                    server.SearchBase,
			"use_base_con":                   server.UseBaseCon,
			"use_subtree":                    server.UseSubtree,
			"use_trust_store":                server.UseTrustStore,
			"user_email_address_attribute":   server.UserEmailAddressAttribute,
			"user_id_attribute":              server.UserIDAttribute,
			"user_ldap_filter":               server.UserLDAPFilter,
			"user_member_of_attribute":       server.UserMemberOfAttribute,
			"user_object_class":              server.UserObjectClass,
			"user_password_attribute":        server.UserPasswordAttribute,
			"user_real_name_attribute":       server.UserRealNameAttribute,
		}
	}
	return nil
}
