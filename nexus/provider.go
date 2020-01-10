package nexus

import (
	"github.com/datadrivers/terraform-provider-nexus/nexus"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider returns a terraform.ResourceProvider
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"nexus_repository": resourceRepository(),
			"nexus_user":       resourceUser(),
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
	config := nexus.Config{
		URL:      d.Get("url").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}
	return nexus.NewClient(config), nil
}
