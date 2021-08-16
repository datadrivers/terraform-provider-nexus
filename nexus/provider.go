package nexus

import (
	client "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"nexus_anonymous":       dataSourceAnonymous(),
			"nexus_blobstore":       dataSourceBlobstore(),
			"nexus_privileges":      dataSourcePrivileges(),
			"nexus_repository":      dataSourceRepository(),
			"nexus_routing_rule":    dataSourceRoutingRule(),
			"nexus_security_ldap":   dataSourceSecurityLDAP(),
			"nexus_security_realms": dataSourceSecurityRealms(),
			"nexus_security_saml":   dataSourceSecuritySAML(),
			"nexus_security_user":   dataSourceSecurityUser(),
			"nexus_user":            dataSourceUser(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"nexus_anonymous":           resourceAnonymous(),
			"nexus_blobstore":           resourceBlobstore(),
			"nexus_content_selector":    resourceContentSelector(),
			"nexus_privilege":           resourcePrivilege(),
			"nexus_repository":          resourceRepository(),
			"nexus_role":                resourceRole(),
			"nexus_routing_rule":        resourceRoutingRule(),
			"nexus_script":              resourceScript(),
			"nexus_security_ldap":       resourceSecurityLDAP(),
			"nexus_security_ldap_order": resourceSecurityLDAPOrder(),
			"nexus_security_realms":     resourceSecurityRealms(),
			"nexus_security_saml":       resourceSecuritySAML(),
			"nexus_security_user":       resourceSecurityUser(),
			"nexus_user":                resourceUser(),
		},
		Schema: map[string]*schema.Schema{
			"insecure": {
				Default:     false,
				DefaultFunc: schema.EnvDefaultFunc("NEXUS_INSECURE_SKIP_VERIFY", "true"),
				Optional:    true,
				Type:        schema.TypeBool,
			},
			"password": {
				DefaultFunc: schema.EnvDefaultFunc("NEXUS_PASSWORD", "admin123"),
				Required:    true,
				Type:        schema.TypeString,
			},
			"url": {
				DefaultFunc: schema.EnvDefaultFunc("NEXUS_URL", "http://127.0.0.1:8080"),
				Required:    true,
				Type:        schema.TypeString,
			},
			"username": {
				DefaultFunc: schema.EnvDefaultFunc("NEXUS_USERNAME", "admin"),
				Required:    true,
				Type:        schema.TypeString,
			},
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := client.Config{
		Insecure: d.Get("insecure").(bool),
		Password: d.Get("password").(string),
		URL:      d.Get("url").(string),
		Username: d.Get("username").(string),
	}
	return client.NewClient(config), nil
}
