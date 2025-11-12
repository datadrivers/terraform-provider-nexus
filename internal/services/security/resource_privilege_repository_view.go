package security

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	nexus "github.com/gcroucher/go-nexus-client/nexus3"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceSecurityPrivilegeRepositoryView() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage a privilege for a repository view",

		Create: resourceSecurityPrivilegeRepositoryViewCreate,
		Read:   resourceSecurityPrivilegeRepositoryViewRead,
		Update: resourceSecurityPrivilegeRepositoryViewUpdate,
		Delete: resourceSecurityPrivilegeRepositoryViewDelete,
		Exists: resourceSecurityPrivilegeRepositoryViewExists,
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
			"repository": {
				Description: "Name of the repository the privilege applies to",
				Type:        schema.TypeString,
				Required:    true,
			},
			"format": {
				Description: "The format of the referenced Repository",
				Type:        schema.TypeString,
				Required:    true,
			},
			"actions": {
				Description: "A list of allowed actions. For a list of applicable values see https://help.sonatype.com/repomanager3/nexus-repository-administration/access-control/privileges#Privileges-PrivilegeTypes",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"ADD", "READ", "DELETE", "BROWSE", "EDIT"}, false),
				},
			},
		},
	}
}

func resourceSecurityPrivilegeRepositoryViewRead(resourceData *schema.ResourceData, m interface{}) error {
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
	resourceData.Set("actions", privilege.Actions)
	resourceData.Set("repository", privilege.Repository)
	resourceData.Set("format", privilege.Format)
	return nil
}

func getRepositoryViewFromResourceData(d *schema.ResourceData) security.PrivilegeRepositoryView {
	var actions []security.SecurityPrivilegeRepositoryViewActions

	var actionsUnwrapped = d.Get("actions").([]interface{})

	for _, element := range actionsUnwrapped {
		actions = append(actions, security.SecurityPrivilegeRepositoryViewActions(element.(string)))
	}

	privilege := security.PrivilegeRepositoryView{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Repository:  d.Get("repository").(string),
		Format:      d.Get("format").(string),
		Actions:     actions,
	}
	return privilege
}

func resourceSecurityPrivilegeRepositoryViewCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	privilege := getRepositoryViewFromResourceData(d)

	if err := client.Security.Privilege.RepositoryView.Create(privilege); err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))
	return resourceSecurityPrivilegeRepositoryViewRead(d, m)
}

func resourceSecurityPrivilegeRepositoryViewUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	privilege := getRepositoryViewFromResourceData(d)
	if err := client.Security.Privilege.RepositoryView.Update(privilege.Name, privilege); err != nil {
		return err
	}

	return resourceSecurityPrivilegeRepositoryViewRead(d, m)
}

func resourceSecurityPrivilegeRepositoryViewDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.Security.Privilege.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceSecurityPrivilegeRepositoryViewExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	user, err := client.Security.Privilege.Get(d.Id())
	return user != nil, err
}
