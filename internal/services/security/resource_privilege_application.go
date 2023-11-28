package security

import (
	nexus "github.com/dre2004/go-nexus-client/nexus3"
	"github.com/dre2004/go-nexus-client/nexus3/schema/security"
	"github.com/dre2004/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceSecurityPrivilegeApplication() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage a privilege for an application",

		Create: resourceSecurityPrivilegeApplicationCreate,
		Read:   resourceSecurityPrivilegeApplicationRead,
		Update: resourceSecurityPrivilegeApplicationUpdate,
		Delete: resourceSecurityPrivilegeApplicationDelete,
		Exists: resourceSecurityPrivilegeApplicationExists,
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
			"domain": {
				Description: "A domain",
				Type:        schema.TypeString,
				Required:    true,
			},
			"actions": {
				Description: "A list of allowed actions. For a list of applicable values see https://help.sonatype.com/repomanager3/nexus-repository-administration/access-control/privileges#Privileges-PrivilegeTypes",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"ADD", "READ", "EDIT", "DELETE"}, false),
				},
			},
		},
	}
}

func resourceSecurityPrivilegeApplicationRead(resourceData *schema.ResourceData, m interface{}) error {
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
	resourceData.Set("domain", privilege.Domain)
	resourceData.Set("actions", privilege.Actions)

	return nil
}

func getPrivilegeApplicationFromResourceData(d *schema.ResourceData) security.PrivilegeApplication {
	var actions []security.SecurityPrivilegeApplicationActions

	var actionsUnwrapped = d.Get("actions").([]interface{})

	for _, element := range actionsUnwrapped {
		actions = append(actions, security.SecurityPrivilegeApplicationActions(element.(string)))
	}

	privilege := security.PrivilegeApplication{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Domain:      d.Get("domain").(string),
		Actions:     actions,
	}
	return privilege
}

func resourceSecurityPrivilegeApplicationCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	privilege := getPrivilegeApplicationFromResourceData(d)

	if err := client.Security.Privilege.Application.Create(privilege); err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))
	return resourceSecurityPrivilegeApplicationRead(d, m)
}

func resourceSecurityPrivilegeApplicationUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	privilege := getPrivilegeApplicationFromResourceData(d)
	if err := client.Security.Privilege.Application.Update(privilege.Name, privilege); err != nil {
		return err
	}

	return resourceSecurityPrivilegeApplicationRead(d, m)
}

func resourceSecurityPrivilegeApplicationDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.Security.Privilege.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceSecurityPrivilegeApplicationExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	privilege, err := client.Security.Privilege.Get(d.Id())
	return privilege != nil, err
}
