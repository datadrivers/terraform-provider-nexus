package nexus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/datadrivers/terraform-provider-nexus/nexus/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	nexusUsersAPIEndpoint = "service/rest/beta/security/users"
)

type NexusUser struct {
	UserID       string   `json:"userId"`
	FirstName    string   `json:"firstName"`
	LastName     string   `json:"lastName"`
	EmailAddress string   `json:"emailAddress"`
	Password     string   `json:"password"`
	Status       string   `json:"status"`
	Roles        []string `json:"roles"`
}

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

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*client.Config)
	client := client.NewClient(*config)

	user := NexusUser{
		UserID:       d.Get("userid").(string),
		FirstName:    d.Get("firstname").(string),
		LastName:     d.Get("lastname").(string),
		EmailAddress: d.Get("email").(string),
		Password:     d.Get("password").(string),
		Status:       d.Get("status").(string),
		Roles:        resourceDataStringSlice(d, "roles"),
	}

	reqBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("could not marshal user data: %v", err)
	}

	reqReader := bytes.NewReader(reqBytes)

	body, resp, err := client.Post(nexusUsersAPIEndpoint, reqReader)
	if err != nil {
		return fmt.Errorf("error while creating user: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error while creating user: HTTP: %d, %s", resp.StatusCode, string(body))
	}

	d.SetId(user.UserID)
	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*client.Config)
	client := client.NewClient(*config)

	userId := d.Get("userid").(string)

	resp, err := client.Delete(fmt.Sprintf("%s/%s", nexusUsersAPIEndpoint, userId))
	if err != nil {
		return fmt.Errorf("could not delete user: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != 204 {
		return fmt.Errorf("could not delete user: HTTP: %d", resp.StatusCode)
	}

	d.SetId("")
	return nil
}

func resourceDataStringSlice(d *schema.ResourceData, attribute string) []string {
	n := d.Get(fmt.Sprintf("%s.#", attribute)).(int)
	data := make([]string, n)
	for i := 0; i < n; i++ {
		data[i] = d.Get(fmt.Sprintf("%s.%s", attribute, i)).(string)
	}
	return data
}
