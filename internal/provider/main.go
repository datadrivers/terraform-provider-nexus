package provider

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/pkg/client"
	"github.com/datadrivers/terraform-provider-nexus/internal/services/blobstore"
	"github.com/datadrivers/terraform-provider-nexus/internal/services/deprecated"
	"github.com/datadrivers/terraform-provider-nexus/internal/services/other"
	"github.com/datadrivers/terraform-provider-nexus/internal/services/repository"
	"github.com/datadrivers/terraform-provider-nexus/internal/services/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a terraform.Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"nexus_anonymous":                 deprecated.DataSourceAnonymous(),
			"nexus_blobstore":                 deprecated.DataSourceBlobstore(),
			"nexus_blobstore_azure":           blobstore.DataSourceBlobstoreAzure(),
			"nexus_blobstore_file":            blobstore.DataSourceBlobstoreFile(),
			"nexus_blobstore_group":           blobstore.DataSourceBlobstoreGroup(),
			"nexus_blobstore_s3":              blobstore.DataSourceBlobstoreS3(),
			"nexus_privileges":                deprecated.DataSourcePrivileges(),
			"nexus_repository":                deprecated.DataSourceRepository(),
			"nexus_repository_apt_hosted":     repository.DataSourceRepositoryAptHosted(),
			"nexus_repository_list":           repository.DataSourceRepositoryList(),
			"nexus_repository_yum_hosted":     repository.DataSourceRepositoryYumHosted(),
			"nexus_routing_rule":              other.DataSourceRoutingRule(),
			"nexus_security_anonymous":        security.DataSourceSecurityAnonymous(),
			"nexus_security_content_selector": security.DataSourceSecurityContentSelector(),
			"nexus_security_ldap":             security.DataSourceSecurityLDAP(),
			"nexus_security_realms":           security.DataSourceSecurityRealms(),
			"nexus_security_role":             security.DataSourceSecurityRole(),
			"nexus_security_saml":             security.DataSourceSecuritySAML(),
			"nexus_security_user":             security.DataSourceSecurityUser(),
			"nexus_security_user_token":       security.DataSourceSecurityUserToken(),
			"nexus_user":                      deprecated.DataSourceUser(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"nexus_anonymous":                 deprecated.ResourceAnonymous(),
			"nexus_blobstore":                 deprecated.ResourceBlobstore(),
			"nexus_blobstore_azure":           blobstore.ResourceBlobstoreAzure(),
			"nexus_blobstore_file":            blobstore.ResourceBlobstoreFile(),
			"nexus_blobstore_group":           blobstore.ResourceBlobstoreGroup(),
			"nexus_blobstore_s3":              blobstore.ResourceBlobstoreS3(),
			"nexus_content_selector":          deprecated.ResourceContentSelector(),
			"nexus_privilege":                 deprecated.ResourcePrivilege(),
			"nexus_repository":                deprecated.ResourceRepository(),
			"nexus_repository_apt_hosted":     repository.ResourceRepositoryAptHosted(),
			"nexus_repository_yum_hosted":     repository.ResourceRepositoryYumHosted(),
			"nexus_role":                      deprecated.ResourceRole(),
			"nexus_routing_rule":              other.ResourceRoutingRule(),
			"nexus_script":                    other.ResourceScript(),
			"nexus_security_anonymous":        security.ResourceSecurityAnonymous(),
			"nexus_security_content_selector": security.ResourceSecurityContentSelector(),
			"nexus_security_ldap":             security.ResourceSecurityLDAP(),
			"nexus_security_ldap_order":       security.ResourceSecurityLDAPOrder(),
			"nexus_security_realms":           security.ResourceSecurityRealms(),
			"nexus_security_role":             security.ResourceSecurityRole(),
			"nexus_security_saml":             security.ResourceSecuritySAML(),
			"nexus_security_user":             security.ResourceSecurityUser(),
			"nexus_security_user_token":       security.ResourceSecurityUserToken(),
			"nexus_user":                      deprecated.ResourceUser(),
		},
		Schema: map[string]*schema.Schema{
			"insecure": {
				Description: "Boolean to specify wether insecure SSL connections are allowed or not. Reading environment variable NEXUS_INSECURE_SKIP_VERIFY. Default:`true`",
				Default:     false,
				DefaultFunc: schema.EnvDefaultFunc("NEXUS_INSECURE_SKIP_VERIFY", "true"),
				Optional:    true,
				Type:        schema.TypeBool,
			},
			"password": {
				Description: "Password of user to connect to API. Reading environment variable NEXUS_PASSWORD. Default:`admin123`",
				DefaultFunc: schema.EnvDefaultFunc("NEXUS_PASSWORD", "admin123"),
				Required:    true,
				Type:        schema.TypeString,
			},
			"url": {
				Description: "URL of Nexus to reach API. Reading environment variable NEXUS_URL. Default:`http://127.0.0.1:8080`",
				DefaultFunc: schema.EnvDefaultFunc("NEXUS_URL", "http://127.0.0.1:8080"),
				Required:    true,
				Type:        schema.TypeString,
			},
			"username": {
				Description: "Username used to connect to API. Reading environment variable NEXUS_USERNAME. Default:`admin`",
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

	return nexus.NewClient(config), nil
}
