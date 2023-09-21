package security

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceSecurityPrivilegeRepositoryContentSelector() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage a privilege for a repository content selector",

		Create: resourceSecurityPrivilegeRepositoryContentSelectorCreate,
		Read:   resourceSecurityPrivilegeRepositoryContentSelectorRead,
		Update: resourceSecurityPrivilegeRepositoryContentSelectorUpdate,
		Delete: resourceSecurityPrivilegeRepositoryContentSelectorDelete,
		Exists: resourceSecurityPrivilegeRepositoryContentSelectorExists,
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
			"content_selector": {
				Description: "The content selector",
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

func resourceSecurityPrivilegeRepositoryContentSelectorRead(resourceData *schema.ResourceData, m interface{}) error {
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
	resourceData.Set("content_selector", privilege.ContentSelector)

	return nil
}

func getPrivilegeRepositoryContentSelectorFromResourceData(d *schema.ResourceData) security.PrivilegeRepositoryContentSelector {
	var actions []security.SecurityPrivilegeRepositoryContentSelectorActions

	var actionsUnwrapped = d.Get("actions").([]interface{})

	for _, element := range actionsUnwrapped {
		actions = append(actions, security.SecurityPrivilegeRepositoryContentSelectorActions(element.(string)))
	}

	privilege := security.PrivilegeRepositoryContentSelector{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		Repository:      d.Get("repository").(string),
		Format:          d.Get("format").(string),
		ContentSelector: d.Get("content_selector").(string),
		Actions:         actions,
	}
	return privilege
}

func resourceSecurityPrivilegeRepositoryContentSelectorCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	privilege := getPrivilegeRepositoryContentSelectorFromResourceData(d)

	if err := client.Security.Privilege.RepositoryContentSelector.Create(privilege); err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))
	return resourceSecurityPrivilegeRepositoryContentSelectorRead(d, m)
}

func resourceSecurityPrivilegeRepositoryContentSelectorUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	privilege := getPrivilegeRepositoryContentSelectorFromResourceData(d)
	if err := client.Security.Privilege.RepositoryContentSelector.Update(privilege.Name, privilege); err != nil {
		return err
	}

	return resourceSecurityPrivilegeRepositoryContentSelectorRead(d, m)
}

func resourceSecurityPrivilegeRepositoryContentSelectorDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.Security.Privilege.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceSecurityPrivilegeRepositoryContentSelectorExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	user, err := client.Security.Privilege.Get(d.Id())
	return user != nil, err
}
