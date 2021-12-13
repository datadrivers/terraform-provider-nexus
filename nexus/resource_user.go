/*
This resource is deprecated. Please use the data source "nexus_security_user" instead.

Use this resource to manage users

Example Usage

```hcl
resource "nexus_user" "admin" {
  userid    = "admin"
  firstname = "Administrator"
  lastname  = "User"
  email     = "nexus@example.com"
  password  = "admin123"
  roles     = ["nx-admin"]
  status    = "active"
}
```
*/
package nexus

import (
	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Exists: resourceUserExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"userid": {
				Description: "The userid which is required for login. This value cannot be changed.",
				ForceNew:    true,
				Type:        schema.TypeString,
				Required:    true,
			},
			"firstname": {
				Description: "The first name of the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"lastname": {
				Description: "The last name of the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"email": {
				Description: "The email address associated with the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": {
				Description: "The password for the user.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"roles": {
				Description: "The roles which the user has been assigned within Nexus.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Type:        schema.TypeSet,
			},
			"status": {
				Default:     "active",
				Description: "The user's status, e.g. active or disabled.",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"active",
					"disabled",
				}, false),
			},
		},
	}
}

func getUserFromResourceData(d *schema.ResourceData) nexus.User {
	return nexus.User{
		UserID:       d.Get("userid").(string),
		FirstName:    d.Get("firstname").(string),
		LastName:     d.Get("lastname").(string),
		EmailAddress: d.Get("email").(string),
		Password:     d.Get("password").(string),
		Status:       d.Get("status").(string),
		Roles:        interfaceSliceToStringSlice(d.Get("roles").(*schema.Set).List()),
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)
	user := getUserFromResourceData(d)

	if err := nexusClient.UserCreate(user); err != nil {
		return err
	}

	d.SetId(user.UserID)
	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	user, err := nexusClient.UserRead(d.Id())
	if err != nil {
		return err
	}

	if user == nil {
		d.SetId("")
		return nil
	}

	d.Set("email", user.EmailAddress)
	d.Set("firstname", user.FirstName)
	d.Set("lastname", user.LastName)
	d.Set("roles", stringSliceToInterfaceSlice(user.Roles))
	d.Set("status", user.Status)
	d.Set("userid", user.UserID)

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	d.Partial(true)

	if d.HasChange("password") {
		password := d.Get("password").(string)
		if err := nexusClient.UserChangePassword(d.Id(), password); err != nil {
			return err
		}
		d.SetPartial("password")
	}

	d.Partial(false)

	if d.HasChange("firstname") || d.HasChange("lastname") || d.HasChange("email") || d.HasChange("status") || d.HasChange("roles") {
		user := getUserFromResourceData(d)
		if err := nexusClient.UserUpdate(d.Id(), user); err != nil {
			return err
		}
	}
	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	if err := nexusClient.UserDelete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceUserExists(d *schema.ResourceData, m interface{}) (bool, error) {
	nexusClient := m.(nexus.Client)

	user, err := nexusClient.UserRead(d.Id())
	return user != nil, err
}
