package nexus

import (
	nexus "github.com/datadrivers/terraform-provider-nexus/nexus/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			"userid": {
				Description: "The userid which is required for login. This value cannot be changed.",
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
			},
			"roles": {
				Description: "The roles which the user has been assigned within Nexus.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Type:        schema.TypeList,
			},
			"status": {
				Description: "The user's status, e.g. active or disabled.",
				Type:        schema.TypeString,
				Required:    true,
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
		Roles:        resourceDataStringSlice(d, "roles"),
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
	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	nexus := m.(nexus.Client)
	if d.HasChange("firstname") || d.HasChange("lastname") || d.HasChange("email") || d.HasChange("status") {
		userId := d.Get("userid").(string)
		user := getUserFromResourceData(d)
		if err := nexus.UserUpdate(userId, user); err != nil {
			return err
		}
	}
	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	nexus := m.(nexus.Client)

	userId := d.Get("userid").(string)

	if err := nexus.UserDelete(userId); err != nil {
		return err
	}

	d.SetId("")
	return nil
}
