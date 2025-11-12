package security

import (
	"strings"

	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	nexus "github.com/gcroucher/go-nexus-client/nexus3"
	"github.com/gcroucher/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSecurityRole() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a Nexus Role.",

		Create: resourceSecurityRoleCreate,
		Read:   resourceSecurityRoleRead,
		Update: resourceSecurityRoleUpdate,
		Delete: resourceSecurityRoleDelete,
		Exists: resourceSecurityRoleExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
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

func getSecurityRoleFromResourceData(d *schema.ResourceData) security.Role {
	return security.Role{
		ID:          d.Get("roleid").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Privileges:  tools.InterfaceSliceToStringSlice(d.Get("privileges").(*schema.Set).List()),
		Roles:       tools.InterfaceSliceToStringSlice(d.Get("roles").(*schema.Set).List()),
	}
}

func resourceSecurityRoleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	role := getSecurityRoleFromResourceData(d)
	if err := client.Security.Role.Create(role); err != nil {
		return err
	}

	d.SetId(role.ID)
	return resourceSecurityRoleRead(d, m)
}

func resourceSecurityRoleRead(d *schema.ResourceData, m interface{}) error {
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

func resourceSecurityRoleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	roleID := d.Get("roleid").(string)

	role := getSecurityRoleFromResourceData(d)
	if err := client.Security.Role.Update(roleID, role); err != nil {
		return err
	}

	return resourceSecurityRoleRead(d, m)
}

func resourceSecurityRoleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.Security.Role.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceSecurityRoleExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	role, err := client.Security.Role.Get(d.Id())
	return role != nil, err
}
