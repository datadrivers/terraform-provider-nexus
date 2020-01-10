package nexus

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// "github.com/sonatype-nexus-community/gonexus/rm"
)

func resourceRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryCreate,
		Read:   resourceRepositoryRead,
		Update: resourceRepositoryUpdate,
		Delete: resourceRepositoryDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "A unique identifier for this repository",
				Type:        schema.TypeString,
				Required:    true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"format": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRepositoryCreate(d *schema.ResourceData, m interface{}) error {
	// config := m.(*Config)

	// rm,err := nexusrm.New(config.URL, config.Username, config.Password)
	// if err != nil {
	// 	return err
	// }
	repoName := d.Get("name").(string)
	// repo := nexusrm.Repository{
	// 	Name: repoName,
	// 	Type: d.Get("type").(string),
	// 	Format: d.Get("format").(string),
	// }

	d.SetId(repoName)
	return resourceRepositoryRead(d, m)
}

func resourceRepositoryRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceRepositoryUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceRepositoryRead(d, m)
}

func resourceRepositoryDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
