package security

import (
	"errors"
	"fmt"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceSecurityUser() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage users.",

		Create: resourceSecurityUserCreate,
		Read:   resourceSecurityUserRead,
		Update: resourceSecurityUserUpdate,
		Delete: resourceSecurityUserDelete,
		Exists: resourceSecurityUserExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
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
				Description:   "The password for the user.",
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"password_wo", "password_wo_version"},
			},
			"password_wo": {
				Description:   "The password for the user (write-only, not stored in state). Use with password_wo_version to control updates.",
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				WriteOnly:     true,
				ConflictsWith: []string{"password"},
				RequiredWith:  []string{"password_wo_version"},
			},
			"password_wo_version": {
				Description:   "Version tracker for password_wo changes. Increment this value to force a password update when using password_wo.",
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"password"},
				RequiredWith:  []string{"password_wo"},
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
			"source": {
				Default:     "default",
				Description: "The user's source, e.g. default (local)",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"default",
				}, false),
			},
		},
	}
}

func getSecurityUserFromResourceData(d *schema.ResourceData) security.User {
	return security.User{
		UserID:       d.Get("userid").(string),
		FirstName:    d.Get("firstname").(string),
		LastName:     d.Get("lastname").(string),
		EmailAddress: d.Get("email").(string),
		Password:     getPasswordFromResourceData(d),
		Status:       d.Get("status").(string),
		Roles:        tools.InterfaceSliceToStringSlice(d.Get("roles").(*schema.Set).List()),
		Source:       d.Get("source").(string),
	}
}

func getPasswordFromResourceData(d *schema.ResourceData) string {
	if password := d.Get("password").(string); password != "" {
		return password
	}

	if d.Get("password_wo_version").(int) != 0 && d.HasChange("password_wo_version") {
		passwordWriteOnly, diags := d.GetRawConfigAt(cty.GetAttrPath("password_wo"))
		if diags.HasError() {
			return ""
		}

		if passwordWriteOnly.IsNull() || !passwordWriteOnly.IsKnown() {
			return ""
		}

		return passwordWriteOnly.AsString()
	}

	return ""
}

func resourceSecurityUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	user := getSecurityUserFromResourceData(d)

	if user.Password == "" {
		return fmt.Errorf("either 'password' or 'password_wo' with 'password_wo_version' must be provided")
	}

	if err := client.Security.User.Create(user); err != nil {
		return err
	}

	d.SetId(user.UserID)
	return resourceSecurityUserRead(d, m)
}

func resourceSecurityUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	user, err := client.Security.User.Get(d.Id(), nil)
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
	d.Set("roles", tools.StringSliceToInterfaceSlice(user.Roles))
	d.Set("status", user.Status)
	d.Set("userid", user.UserID)
	d.Set("source", user.Source)

	if v, ok := d.GetOk("password_wo_version"); ok && v != nil {
		d.Set("password_wo_version", v.(int))
	} else {
		if existingPassword, hasPassword := d.GetOk("password"); hasPassword {
			d.Set("password", existingPassword.(string))
		}
	}

	return nil
}

func resourceSecurityUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	passwordChanged := false
	var newPassword string

	if d.HasChange("password") {
		newPassword = d.Get("password").(string)
		passwordChanged = true
	}

	if d.HasChange("password_wo_version") {
		passwordWriteOnly, diags := d.GetRawConfigAt(cty.GetAttrPath("password_wo"))
		if diags.HasError() {
			return errors.New("error reading 'password_wo' argument")
		}

		if !passwordWriteOnly.IsNull() && passwordWriteOnly.IsKnown() {
			newPassword = passwordWriteOnly.AsString()
			passwordChanged = true
		}
	}

	if passwordChanged && newPassword != "" {
		if err := client.Security.User.ChangePassword(d.Id(), newPassword); err != nil {
			return err
		}
	}

	if d.HasChange("firstname") || d.HasChange("lastname") || d.HasChange("email") || d.HasChange("status") || d.HasChange("roles") {
		user := getSecurityUserFromResourceData(d)
		if !passwordChanged {
			user.Password = ""
		}
		if err := client.Security.User.Update(d.Id(), user); err != nil {
			return err
		}
	}

	return resourceSecurityUserRead(d, m)
}

func resourceSecurityUserDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.Security.User.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceSecurityUserExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	user, err := client.Security.User.Get(d.Id(), nil)
	return user != nil, err
}
