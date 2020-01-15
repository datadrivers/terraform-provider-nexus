package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryCreate,
		Read:   resourceRepositoryRead,
		Update: resourceRepositoryUpdate,
		Delete: resourceRepositoryDelete,

		Schema: map[string]*schema.Schema{
			"format": {
				Required: true,
				Type:     schema.TypeString,
			},
			"name": {
				Description: "A unique identifier for this repository",
				Required:    true,
				Type:        schema.TypeString,
			},
			"online": {
				Default:     true,
				Description: "",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func getRepositoryFromResourceData(d *schema.ResourceData) nexus.Repository {
	return nexus.Repository{
		Name:   d.Get("name").(string),
		Online: d.Get("online").(bool),
	}
}

func resourceRepositoryCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(nexus.Client)
	repo := getRepositoryFromResourceData(d)
	repoFormat := d.Get("format").(string)
	repoType := d.Get("type").(string)

	if err := client.RepositoryCreate(repo, repoFormat, repoType); err != nil {
		return err
	}

	d.SetId(repo.Name)
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
