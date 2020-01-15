package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceRoleCreate,
		Read:   resourceRoleRead,
		Update: resourceRoleUpdate,
		Delete: resourceRoleDelete,

		Schema: map[string]*schema.Schema{
			"roleid": {
				Description: "The id of the role.",
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
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Type:        schema.TypeList,
			},
			"roles": {
				Description: "The roles of this role.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Type:        schema.TypeList,
			},
		},
	}
}

func getRoleFromResourceData(d *schema.ResourceData) nexus.Role {
	return nexus.Role{
		ID:          d.Get("roleid").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Privileges:  resourceDataStringSlice(d, "privileges"),
		Roles:       resourceDataStringSlice(d, "roles"),
	}
}

func resourceRoleCreate(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)
	role := getRoleFromResourceData(d)

	if err := nexusClient.RoleCreate(role); err != nil {
		return err
	}

	d.SetId(role.ID)
	return resourceRoleRead(d, m)
}

func resourceRoleRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceRoleUpdate(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)
	roleId := d.Get("roleid").(string)

	if d.HasChange("name") || d.HasChange("description") || d.HasChange("privileges") || d.HasChange("roles") {
		role := getRoleFromResourceData(d)
		if err := nexusClient.RoleUpdate(roleId, role); err != nil {
			return err
		}
	}

	return resourceRoleRead(d, m)
}

func resourceRoleDelete(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	roleId := d.Get("roleid").(string)

	if err := nexusClient.RoleDelete(roleId); err != nil {
		return err
	}

	d.SetId("")
	return nil
}
