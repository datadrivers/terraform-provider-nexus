package security

import (
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	nexus "github.com/gcroucher/go-nexus-client/nexus3"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceSecurityPrivilegeRepositoryAdmin() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage a privilege for a repository admin",

		Create: resourceSecurityPrivilegeRepositoryAdminCreate,
		Read:   resourceSecurityPrivilegeRepositoryAdminRead,
		Update: resourceSecurityPrivilegeRepositoryAdminUpdate,
		Delete: resourceSecurityPrivilegeRepositoryAdminDelete,
		Exists: resourceSecurityPrivilegeRepositoryAdminExists,
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

func resourceSecurityPrivilegeRepositoryAdminRead(resourceData *schema.ResourceData, m interface{}) error {
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

func getPrivilegeRepositoryAdminFromResourceData(d *schema.ResourceData) security.PrivilegeRepositoryAdmin {
	var actions []security.SecurityPrivilegeRepositoryAdminActions

	var actionsUnwrapped = d.Get("actions").([]interface{})

	for _, element := range actionsUnwrapped {
		actions = append(actions, security.SecurityPrivilegeRepositoryAdminActions(element.(string)))
	}

	privilege := security.PrivilegeRepositoryAdmin{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Repository:  d.Get("repository").(string),
		Format:      d.Get("format").(string),
		Actions:     actions,
	}
	return privilege
}

func resourceSecurityPrivilegeRepositoryAdminCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	privilege := getPrivilegeRepositoryAdminFromResourceData(d)

	if err := client.Security.Privilege.RepositoryAdmin.Create(privilege); err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))
	return resourceSecurityPrivilegeRepositoryAdminRead(d, m)
}

func resourceSecurityPrivilegeRepositoryAdminUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	privilege := getPrivilegeRepositoryAdminFromResourceData(d)
	if err := client.Security.Privilege.RepositoryAdmin.Update(privilege.Name, privilege); err != nil {
		return err
	}

	return resourceSecurityPrivilegeRepositoryAdminRead(d, m)
}

func resourceSecurityPrivilegeRepositoryAdminDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.Security.Privilege.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceSecurityPrivilegeRepositoryAdminExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	user, err := client.Security.Privilege.Get(d.Id())
	return user != nil, err
}
