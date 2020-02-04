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
			"nexus_blobstore":  dataSourceBlobstore(),
			"nexus_repository": dataSourceRepository(),
			"nexus_user":       dataSourceUser(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"nexus_repository": resourceRepository(),
			"nexus_role":       resourceRole(),
			"nexus_user":       resourceUser(),
			"nexus_script":     resourceScript(),
		},
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEXUS_URL", "http://127.0.0.1:8080"),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEXUS_USERNAME", "admin"),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEXUS_PASSWORD", "admin123"),
			},
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := client.Config{
		URL:      d.Get("url").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}
	return client.NewClient(config), nil
}
