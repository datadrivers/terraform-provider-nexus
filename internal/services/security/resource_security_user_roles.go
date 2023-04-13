package security

import (
	"errors"
	"fmt"
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/terraform-provider-nexus/internal/schema/common"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"sort"
)

func ResourceSecurityUserRoles() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage roles of existing users.\n\n" +
			"Conflicts with nexus_security_user, it will cause drifts",

		Create: resourceSecurityUserRolesCreate,
		Read:   resourceSecurityUserRolesRead,
		Update: resourceSecurityUserRolesUpdate,
		Delete: resourceSecurityUserRolesDelete,
		Exists: resourceSecurityUserRolesExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": common.ResourceID,
			"userid": {
				Description: "User ID of the user whose roles you wish to manage using Terraform. Must exist in Nexus.",
				ForceNew:    true,
				Type:        schema.TypeString,
				Required:    true,
			},
			"roles": {
				Description: "The roles which the user will be assigned within Nexus.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeSet,
				Required:    true,
			},
		},
	}
}

// Function to remove duplicate elements from array
// https://codereview.stackexchange.com/a/192954
func resourceSecurityUserRolesUnique(slice []string) []string {
	// create a map with all the values as key
	uniqMap := make(map[string]struct{})
	for _, v := range slice {
		uniqMap[v] = struct{}{}
	}

	// turn the map keys into a slice
	uniqSlice := make([]string, 0, len(uniqMap))
	for v := range uniqMap {
		uniqSlice = append(uniqSlice, v)
	}
	return uniqSlice
}

// Create a Terraform representation for an existing Nexus user with Terraform-managed roles
// Note: we have to load the whole Nexus user object, patch it and then write it back to Nexus as Nexus doesn't support
// PATCH, only PUT.
func resourceSecurityUserRolesCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[TRACE] resourceSecurityUserRolesCreate called")
	client := m.(*nexus.NexusClient)

	// First, try to load the user
	userId := d.Get("userid").(string)
	log.Printf("[TRACE] Calling User.Get(%s)", d.Id())
	user, err := client.Security.User.Get(userId)
	if err != nil {
		log.Printf("[ERROR] User.Get(%s) failed", d.Id())
		return err
	}
	if user == nil {
		d.SetId("")
		return errors.New(fmt.Sprintf("User.Get(%s) from Nexus is nil", userId))
	}
	log.Printf("[TRACE] Got user object:\n%+v\n", user)

	// Now, set the Roles with the ones we specified - merging with existing roles in the process!
	// Ideally, these should only be {"nx-anonymous"} (every user has at least this)
	// Unfortunately, this *will* lead to drift if nx-anonymous or other roles the user has are not specified in the
	// resource definition.
	newRoles := tools.InterfaceSliceToStringSlice(d.Get("roles").(*schema.Set).List())
	// Create a new slice that holds both the old and new roles as well as sorts and de-duplicates them
	// as Nexus will throw an error otherwise if it encounters duplicates.
	// See https://stackoverflow.com/a/58726780
	finalRoles := make([]string, len(user.Roles), len(user.Roles)+len(newRoles))
	copy(finalRoles, user.Roles)
	finalRoles = append(finalRoles, newRoles...)
	sort.Strings(finalRoles)
	user.Roles = resourceSecurityUserRolesUnique(finalRoles)

	// Finally, write the patched object to Nexus.
	log.Printf("[TRACE] Writing user object:\n%+v\n", user)
	if err := client.Security.User.Update(userId, *user); err != nil {
		log.Printf("[ERROR] User.Update(%s) failed", userId)
		return err
	}
	d.SetId(user.UserID)
	return resourceSecurityUserRead(d, m)
}

// Update the state for a given Nexus user
func resourceSecurityUserRolesRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[TRACE] resourceSecurityUserRolesRead called")
	client := m.(*nexus.NexusClient)

	userId := d.Id()
	log.Printf("[TRACE] Calling User.Get(%s)", d.Id())
	user, err := client.Security.User.Get(userId)
	if err != nil {
		log.Printf("[ERROR] User.Get(%s) failed", d.Id())
		return err
	}
	if user == nil {
		d.SetId("")
		return errors.New(fmt.Sprintf("User.Get(%s) from Nexus is nil", userId))
	}
	log.Printf("[TRACE] Got user object:\n%+v\n", user)

	err = d.Set("roles", tools.StringSliceToInterfaceSlice(user.Roles))
	if err != nil {
		log.Printf("[ERROR] d.Set(roles) failed")
		return err
	}
	err = d.Set("userid", user.UserID)
	if err != nil {
		log.Printf("[ERROR] d.Set(userid) failed")
		return err
	}

	return nil
}

// Update a Nexus user's roles if drift detection recognizes a change
func resourceSecurityUserRolesUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] resourceSecurityUserRolesUpdate called")
	client := m.(*nexus.NexusClient)

	if d.HasChange("roles") {
		userId := d.Id()
		log.Printf("[TRACE] Calling User.Get(%s)", d.Id())
		user, err := client.Security.User.Get(userId)
		if err != nil {
			log.Printf("[ERROR] User.Get(%s) failed", d.Id())
			return err
		}
		if user == nil {
			d.SetId("")
			return errors.New(fmt.Sprintf("User.Get(%s) from Nexus is nil", userId))
		}
		log.Printf("[TRACE] Got user object:\n%+v\n", user)

		finalRoles := tools.InterfaceSliceToStringSlice(d.Get("roles").(*schema.Set).List())
		sort.Strings(finalRoles)
		user.Roles = resourceSecurityUserRolesUnique(finalRoles)

		log.Printf("[DEBUG] Writing user object:\n%+v\n", user)
		if err := client.Security.User.Update(userId, *user); err != nil {
			log.Printf("[ERROR] User.Update(%s) failed", userId)
			return err
		}
	}
	return resourceSecurityUserRolesRead(d, m)
}

// Delete all managed Nexus user roles
// Note: we do NOT delete the user object here because the user may originate from somewhere else (e.g. LDAP)
// Instead we remove all their roles but nx-anonymous
// TODO: store roles that an user had prior to creation in an attribute and restore that?
func resourceSecurityUserRolesDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] resourceSecurityUserRolesDelete called")
	client := m.(*nexus.NexusClient)

	// We only manage roles here. So, we simply set the user's roles to nx-anonymous
	// as we can't pass an empty roles array (the Nexus API will error out
	userId := d.Id()
	log.Printf("[TRACE] Calling User.Get(%s)", d.Id())
	user, err := client.Security.User.Get(userId)
	if err != nil {
		log.Printf("[ERROR] User.Get(%s) failed", d.Id())
		return err
	}
	if user == nil {
		d.SetId("")
		return errors.New(fmt.Sprintf("User.Get(%s) from Nexus is nil", userId))
	}
	log.Printf("[TRACE] Got user object:\n%+v\n", user)

	user.Roles = []string{"nx-anonymous"}

	log.Printf("[DEBUG] Writing user object:\n%+v\n", user)
	if err := client.Security.User.Update(userId, *user); err != nil {
		log.Printf("[ERROR] User.Update(%s) failed", userId)
		return err
	}
	d.SetId("")
	return nil
}

// Check if Nexus has a user
func resourceSecurityUserRolesExists(d *schema.ResourceData, m interface{}) (bool, error) {
	log.Printf("[DEBUG] resourceSecurityUserRolesExists called")
	client := m.(*nexus.NexusClient)

	user, err := client.Security.User.Get(d.Id())
	if err != nil {
		log.Printf("[ERROR] User.Get failed")
	}
	return user != nil, err
}
