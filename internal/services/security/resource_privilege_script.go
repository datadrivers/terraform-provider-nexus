package security

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceSecurityPrivilegeScript() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage a privilege for a Script",

		Create: resourceSecurityPrivilegeScriptCreate,
		Read:   resourceSecurityPrivilegeScriptRead,
		Update: resourceSecurityPrivilegeScriptUpdate,
		Delete: resourceSecurityPrivilegeScriptDelete,
		Exists: resourceSecurityPrivilegeScriptExists,
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
			"script_name": {
				Description: "The script Name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"actions": {
				Description: "A list of allowed actions. For a list of applicable values see https://help.sonatype.com/repomanager3/nexus-repository-administration/access-control/privileges#Privileges-PrivilegeTypes",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"ADD", "READ", "DELETE", "RUN", "BROWSE", "EDIT", "ALL"}, false),
				},
			},
		},
	}
}

func resourceSecurityPrivilegeScriptRead(resourceData *schema.ResourceData, m interface{}) error {
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
	resourceData.Set("script_name", privilege.ScriptName)
	resourceData.Set("actions", privilege.Actions)

	return nil
}

func getPrivilegeScriptFromResourceData(d *schema.ResourceData) security.PrivilegeScript {
	var actions []security.SecurityPrivilegeScriptActions

	var actionsUnwrapped = d.Get("actions").([]interface{})

	for _, element := range actionsUnwrapped {
		actions = append(actions, security.SecurityPrivilegeScriptActions(element.(string)))
	}

	privilege := security.PrivilegeScript{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		ScriptName:  d.Get("script_name").(string),
		Actions:     actions,
	}
	return privilege
}

func resourceSecurityPrivilegeScriptCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	privilege := getPrivilegeScriptFromResourceData(d)

	if err := client.Security.Privilege.Script.Create(privilege); err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))
	return resourceSecurityPrivilegeScriptRead(d, m)
}

func resourceSecurityPrivilegeScriptUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	privilege := getPrivilegeScriptFromResourceData(d)
	if err := client.Security.Privilege.Script.Update(privilege.Name, privilege); err != nil {
		return err
	}

	return resourceSecurityPrivilegeScriptRead(d, m)
}

func resourceSecurityPrivilegeScriptDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.Security.Privilege.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceSecurityPrivilegeScriptExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	user, err := client.Security.Privilege.Get(d.Id())
	return user != nil, err
}
