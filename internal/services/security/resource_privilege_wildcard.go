package security

import (
	nexus "github.com/dre2004/go-nexus-client/nexus3"
	"github.com/dre2004/go-nexus-client/nexus3/schema/security"
	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSecurityPrivilegeWildcard() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage a wildcard privilege",

		Create: resourceSecurityPrivilegeWildcardCreate,
		Read:   resourceSecurityPrivilegeWildcardRead,
		Update: resourceSecurityPrivilegeWildcardUpdate,
		Delete: resourceSecurityPrivilegeWildcardDelete,
		Exists: resourceSecurityPrivilegeWildcardExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
			"name": {
				Description: "The name of the privilege. This value cannot be changed.",
				ForceNew:    true,
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "A description",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"pattern": {
				Description: "The privilege pattern",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceSecurityPrivilegeWildcardRead(resourceData *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	privilege, err := client.Security.Privilege.Get(resourceData.Id())
	if err != nil {
		return err
	}

	if privilege == nil {
		resourceData.SetId("")
		return nil
	}

	resourceData.Set("name", privilege.Name)
	resourceData.Set("description", privilege.Description)
	resourceData.Set("pattern", privilege.Pattern)

	return nil
}

func getPrivilegeWildcardFromResourceData(d *schema.ResourceData) security.PrivilegeWildcard {

	privilege := security.PrivilegeWildcard{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Pattern:     d.Get("pattern").(string),
	}
	return privilege
}

func resourceSecurityPrivilegeWildcardCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	privilege := getPrivilegeWildcardFromResourceData(d)

	if err := client.Security.Privilege.Wildcard.Create(privilege); err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))
	return resourceSecurityPrivilegeWildcardRead(d, m)
}

func resourceSecurityPrivilegeWildcardUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	privilege := getPrivilegeWildcardFromResourceData(d)
	if err := client.Security.Privilege.Wildcard.Update(privilege.Name, privilege); err != nil {
		return err
	}

	return resourceSecurityPrivilegeWildcardRead(d, m)
}

func resourceSecurityPrivilegeWildcardDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.Security.Privilege.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceSecurityPrivilegeWildcardExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	privilege, err := client.Security.Privilege.Get(d.Id())
	return privilege != nil, err
}
