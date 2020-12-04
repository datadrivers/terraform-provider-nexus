/*
Use this resource to create a Nexus Role.

Example Usage

Example Usage - Create a group with roles

```hcl
resource "nexus_role" "nx-admin" {
  roleid      = "nx-admin"
  name        = "nx-admin"
  description = "Administrator role"
  privileges  = ["nx-all"]
  roles       = []
}
```

Example Usage - Create a group with privileges

```hcl
resource "nexus_role" "docker_deploy" {
  description = "Docker deployment role"
  name        = "docker-deploy"
  privileges = [
    "nx-repository-view-docker-*-*",
  ]
  roleid = "docker-deploy"
}
```
*/
package nexus

import (
	"strings"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceRoleCreate,
		Read:   resourceRoleRead,
		Update: resourceRoleUpdate,
		Delete: resourceRoleDelete,
		Exists: resourceRoleExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func getRoleFromResourceData(d *schema.ResourceData) nexus.Role {
	return nexus.Role{
		ID:          d.Get("roleid").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Privileges:  interfaceSliceToStringSlice(d.Get("privileges").(*schema.Set).List()),
		Roles:       interfaceSliceToStringSlice(d.Get("roles").(*schema.Set).List()),
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
	nexusClient := m.(nexus.Client)

	role, err := nexusClient.RoleRead(d.Id())
	if err != nil {
		return err
	}

	if role == nil {
		d.SetId("")
		return nil
	}

	d.Set("description", role.Description)
	d.Set("name", role.Name)
	d.Set("privileges", stringSliceToInterfaceSlice(role.Privileges))
	d.Set("roleid", role.ID)
	d.Set("roles", stringSliceToInterfaceSlice(role.Roles))

	return nil
}

func resourceRoleUpdate(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)
	roleID := d.Get("roleid").(string)

	role := getRoleFromResourceData(d)
	if err := nexusClient.RoleUpdate(roleID, role); err != nil {
		return err
	}

	return resourceRoleRead(d, m)
}

func resourceRoleDelete(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	if err := nexusClient.RoleDelete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceRoleExists(d *schema.ResourceData, m interface{}) (bool, error) {
	nexusClient := m.(nexus.Client)

	role, err := nexusClient.RoleRead(d.Id())
	return role != nil, err
}
