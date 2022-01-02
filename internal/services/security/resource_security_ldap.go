package security

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceSecurityLDAP() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a Nexus Security LDAP configuration.",

		Create: resourceSecurityLDAPCreate,
		Read:   resourceSecurityLDAPRead,
		Update: resourceSecurityLDAPUpdate,
		Delete: resourceSecurityLDAPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"auth_password": {
				Description: "The password to bind with. Required if authScheme other than none.",
				Optional:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
			},
			"auth_realm": {
				Description: "The SASL realm to bind to. Required if authScheme is CRAM_MD5 or DIGEST_MD5",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"auth_schema": {
				Description:  "Authentication scheme used for connecting to LDAP server",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"NONE", "CRAM_MD5", "DIGEST_MD5", "SIMPLE"}, false),
			},
			"auth_username": {
				Description: "This must be a fully qualified username if simple authentication is used. Required if authScheme other than none.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"connection_retry_delay_seconds": {
				Description: "How long to wait before retrying",
				Required:    true,
				Type:        schema.TypeInt,
			},
			"connection_timeout_seconds": {
				Description:  "How long to wait before timeout",
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(1, 3600),
			},
			"group_base_dn": {
				Description: "The relative DN where group objects are found (e.g. ou=Group). This value will have the Search base DN value appended to form the full Group search base DN.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"group_id_attribute": {
				Description: "This field specifies the attribute of the Object class that defines the Group ID. Required if groupType is static",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"group_member_attribute": {
				Description: "LDAP attribute containing the usernames for the group. Required if groupType is static",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"group_member_format": {
				Description: "The format of user ID stored in the group member attribute. Required if groupType is static",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"group_object_class": {
				Description: "LDAP class for group objects. Required if groupType is static",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"group_subtree": {
				Description: "Are groups located in structures below the group base DN",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			"group_type": {
				Description:  "Defines a type of groups used: static (a group contains a list of users) or dynamic (a user contains a list of groups). Required if ldapGroupsAsRoles is true.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"dynamic", "static"}, false),
			},
			"host": {
				Description: "LDAP server connection hostname",
				Required:    true,
				Type:        schema.TypeString,
			},
			"ldap_groups_as_roles": {
				Description: "Denotes whether LDAP assigned roles are used as Nexus Repository Manager roles",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			"max_incident_count": {
				Description: "How many retry attempts",
				Required:    true,
				Type:        schema.TypeInt,
			},
			"name": {
				Description: "LDAP server name",
				Required:    true,
				Type:        schema.TypeString,
			},
			"port": {
				Description: "LDAP server connection port to use",
				Required:    true,
				Type:        schema.TypeInt,
			},
			"protocol": {
				Description:  "LDAP server connection Protocol to use",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"LDAP", "LDAPS"}, true),
			},
			"search_base": {
				Description: "LDAP location to be added to the connection URL",
				Required:    true,
				Type:        schema.TypeString,
			},
			"use_trust_store": {
				Description: "Whether to use certificates stored in Nexus Repository Manager's truststore",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			"user_base_dn": {
				Description: "The relative DN where user objects are found (e.g. ou=people). This value will have the Search base DN value appended to form the full User search base DN.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"user_email_address_attribute": {
				Description: "This is used to find an email address given the user ID",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"user_id_attribute": {
				Description: "This is used to find a user given its user ID",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"user_ldap_filter": {
				Description: "LDAP search filter to limit user search",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"user_member_of_attribute": {
				Description: "Set this to the attribute used to store the attribute which holds groups DN in the user object. Required if groupType is dynamic",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"user_object_class": {
				Description: "LDAP class for user objects",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"user_password_attribute": {
				Description: "If this field is blank the user will be authenticated against a bind with the LDAP server",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"user_real_name_attribute": {
				Description: "This is used to find a real name given the user ID",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"user_subtree": {
				Description: "Are users located in structures below the user base DN?",
				Optional:    true,
				Type:        schema.TypeBool,
			},
		},
	}
}

func resourceSecurityLDAPCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	ldap := getSecurityLDAPFromResourceData(d)

	if err := client.Security.LDAP.Create(ldap); err != nil {
		return err
	}

	if err := setSecurityLDAPToResourceData(&ldap, d); err != nil {
		return err
	}

	return resourceSecurityLDAPRead(d, m)
}

func resourceSecurityLDAPRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	ldap, err := client.Security.LDAP.Get(d.Id())
	if err != nil {
		return err
	}

	if ldap == nil {
		d.SetId("")
		return nil
	}

	return setSecurityLDAPToResourceData(ldap, d)
}

func resourceSecurityLDAPUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	ldapID := d.Id()
	ldap := getSecurityLDAPFromResourceData(d)

	if err := client.Security.LDAP.Update(ldapID, ldap); err != nil {
		return err
	}

	if err := setSecurityLDAPToResourceData(&ldap, d); err != nil {
		return err
	}

	return resourceSecurityLDAPRead(d, m)
}

func resourceSecurityLDAPDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	return client.Security.LDAP.Delete(d.Id())
}

func setSecurityLDAPToResourceData(ldap *security.LDAP, d *schema.ResourceData) error {
	d.SetId(ldap.Name)
	// d.Set("auth_password", ldap.AuthPassword) // AuthPassword is not returned by API
	d.Set("auth_realm", ldap.AuthRealm)
	d.Set("auth_schema", ldap.AuthSchema)
	d.Set("auth_username", ldap.AuthUserName)
	d.Set("connection_retry_delay_seconds", ldap.ConnectionRetryDelaySeconds)
	d.Set("connection_timeout_seconds", ldap.ConnectionTimeoutSeconds)
	d.Set("group_base_dn", ldap.GroupBaseDn)
	d.Set("group_id_attribute", ldap.GroupIDAttribute)
	d.Set("group_member_attribute", ldap.GroupMemberAttribute)
	d.Set("group_member_format", ldap.GroupMemberFormat)
	d.Set("group_object_class", ldap.GroupObjectClass)
	d.Set("group_subtree", ldap.GroupSubtree)
	// d.Set("group_type", ldap.GroupType) // GroupType is not returned by API :-/
	d.Set("host", ldap.Host)
	d.Set("ldap_groups_as_roles", ldap.LDAPGroupsAsRoles)
	d.Set("max_incident_count", ldap.MaxIncidentCount)
	d.Set("name", ldap.Name)
	d.Set("port", ldap.Port)
	d.Set("protocol", ldap.Protocol)
	d.Set("search_base", ldap.SearchBase)
	d.Set("use_trust_store", ldap.UseTrustStore)
	d.Set("user_base_dn", ldap.UserBaseDN)
	d.Set("user_email_address_attribute", ldap.UserEmailAddressAttribute)
	d.Set("user_id_attribute", ldap.UserIDAttribute)
	d.Set("user_ldap_filter", ldap.UserLDAPFilter)
	d.Set("user_member_of_attribute", ldap.UserMemberOfAttribute)
	d.Set("user_object_class", ldap.UserObjectClass)
	d.Set("user_password_attribute", ldap.UserPasswordAttribute)
	d.Set("user_real_name_attribute", ldap.UserRealNameAttribute)
	d.Set("user_subtree", ldap.UserSubtree)

	return nil
}

func getSecurityLDAPFromResourceData(d *schema.ResourceData) security.LDAP {
	ldap := security.LDAP{
		AuthPassword:                d.Get("auth_password").(string),
		AuthRealm:                   d.Get("auth_realm").(string),
		AuthSchema:                  d.Get("auth_schema").(string),
		AuthUserName:                d.Get("auth_username").(string),
		ConnectionRetryDelaySeconds: int32(d.Get("connection_retry_delay_seconds").(int)),
		ConnectionTimeoutSeconds:    int32(d.Get("connection_timeout_seconds").(int)),
		GroupBaseDn:                 d.Get("group_base_dn").(string),
		GroupIDAttribute:            d.Get("group_id_attribute").(string),
		GroupMemberAttribute:        d.Get("group_member_attribute").(string),
		GroupMemberFormat:           d.Get("group_member_format").(string),
		GroupObjectClass:            d.Get("group_object_class").(string),
		GroupSubtree:                d.Get("group_subtree").(bool),
		GroupType:                   d.Get("group_type").(string),
		Host:                        d.Get("host").(string),
		LDAPGroupsAsRoles:           d.Get("ldap_groups_as_roles").(bool),
		MaxIncidentCount:            int32(d.Get("max_incident_count").(int)),
		Name:                        d.Get("name").(string),
		Port:                        int32(d.Get("port").(int)),
		Protocol:                    d.Get("protocol").(string),
		SearchBase:                  d.Get("search_base").(string),
		UseTrustStore:               d.Get("use_trust_store").(bool),
		UserBaseDN:                  d.Get("user_base_dn").(string),
		UserEmailAddressAttribute:   d.Get("user_email_address_attribute").(string),
		UserIDAttribute:             d.Get("user_id_attribute").(string),
		UserLDAPFilter:              d.Get("user_ldap_filter").(string),
		UserMemberOfAttribute:       d.Get("user_member_of_attribute").(string),
		UserObjectClass:             d.Get("user_object_class").(string),
		UserPasswordAttribute:       d.Get("user_password_attribute").(string),
		UserRealNameAttribute:       d.Get("user_real_name_attribute").(string),
		UserSubtree:                 d.Get("user_subtree").(bool),
	}

	return ldap
}
