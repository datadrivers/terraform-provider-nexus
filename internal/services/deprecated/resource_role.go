package deprecated

import (
	"strings"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRole() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource is deprecated. Please use the resource nexus_security_role instead.",
		Description: `!> This resource is deprecated. Please use the resource "nexus_security_role" instead.

Use this resource to create a Nexus Role.`,

		Create: resourceRoleCreate,
		Read:   resourceRoleRead,
		Update: resourceRoleUpdate,
		Delete: resourceRoleDelete,
		Exists: resourceRoleExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"roleid": {
				Description: "The id of the role.",
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeString,
			},
			"name": {
				Description: "The name of the role.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"description": {
				Description: "The description of this role.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"privileges": {
				Description: "The privileges of this role.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Set: func(v interface{}) int {
					return schema.HashString(strings.ToLower(v.(string)))
				},
				Type: schema.TypeSet,
			},
			"roles": {
				Description: "The roles of this role.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Set: func(v interface{}) int {
					return schema.HashString(strings.ToLower(v.(string)))
				},
				Type: schema.TypeSet,
			},
		},
	}
}

func getRoleFromResourceData(d *schema.ResourceData) security.Role {
	return security.Role{
		ID:          d.Get("roleid").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Privileges:  tools.InterfaceSliceToStringSlice(d.Get("privileges").(*schema.Set).List()),
		Roles:       tools.InterfaceSliceToStringSlice(d.Get("roles").(*schema.Set).List()),
	}
}

func resourceRoleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	role := getRoleFromResourceData(d)
	if err := client.Security.Role.Create(role); err != nil {
		return err
	}

	d.SetId(role.ID)
	return resourceRoleRead(d, m)
}

func resourceRoleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	role, err := client.Security.Role.Get(d.Id())
	if err != nil {
		return err
	}

	if role == nil {
		d.SetId("")
		return nil
	}

	d.Set("description", role.Description)
	d.Set("name", role.Name)
	d.Set("privileges", tools.StringSliceToInterfaceSlice(role.Privileges))
	d.Set("roleid", role.ID)
	d.Set("roles", tools.StringSliceToInterfaceSlice(role.Roles))

	return nil
}

func resourceRoleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	roleID := d.Get("roleid").(string)

	role := getRoleFromResourceData(d)
	if err := client.Security.Role.Update(roleID, role); err != nil {
		return err
	}

	return resourceRoleRead(d, m)
}

func resourceRoleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.Security.Role.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceRoleExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	role, err := client.Security.Role.Get(d.Id())
	return role != nil, err
}
