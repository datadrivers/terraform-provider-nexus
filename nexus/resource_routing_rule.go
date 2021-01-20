/*
Use this resource to create a Nexus Routing Rule.

Example Usage

```hcl
resource "nexus_routing_rule" "stop_leaks" {
  name        = "stop-leaks"
  description = "Prevent requests of internal names"
  mode        = "BLOCK"
  matchers    = [
	"^/com/example/.*",
	"^/org/example/.*",
  ]
}
```
*/
package nexus

import (
	"strings"

	nexus "github.com/datadrivers/go-nexus-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceRoutingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceRoutingRuleCreate,
		Read:   resourceRoutingRuleRead,
		Update: resourceRoutingRuleUpdate,
		Delete: resourceRoutingRuleDelete,
		Exists: resourceRoutingRuleExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the routing rule",
				ForceNew:    true,
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the routing rule",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"mode": {
				Description:  "The mode describe how to hande with mathing requests. Possible values: `BLOCK` or `ALLOW` Default: `BLOCK`",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "BLOCK",
				ValidateFunc: validation.StringInSlice([]string{"BLOCK", "ALLOW"}, false),
			},
			"matchers": {
				Description: "Matchers is a list of regular expressions used to identify request paths that are allowed or blocked (depending on above mode)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				Set: func(v interface{}) int {
					return schema.HashString(strings.ToLower(v.(string)))
				},
				Type: schema.TypeSet,
			},
		},
	}
}

func getRoutingRuleFromResourceData(d *schema.ResourceData) nexus.RoutingRule {
	return nexus.RoutingRule{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Mode:        nexus.RoutingRuleMode(d.Get("mode").(string)),
		Matchers:    interfaceSliceToStringSlice(d.Get("matchers").(*schema.Set).List()),
	}
}

func resourceRoutingRuleCreate(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)
	rule := getRoutingRuleFromResourceData(d)

	if err := nexusClient.RoutingRuleCreate(&rule); err != nil {
		return err
	}

	d.SetId(rule.Name)
	return resourceRoutingRuleRead(d, m)
}

func resourceRoutingRuleRead(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	rule, err := nexusClient.RoutingRuleRead(d.Id())
	if err != nil {
		return err
	}

	if rule == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", rule.Name)
	d.Set("description", rule.Description)
	d.Set("mode", rule.Mode)
	d.Set("matchers", stringSliceToInterfaceSlice(rule.Matchers))

	return nil
}

func resourceRoutingRuleUpdate(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	rule := getRoutingRuleFromResourceData(d)
	if err := nexusClient.RoutingRuleUpdate(&rule); err != nil {
		return err
	}

	return resourceRoutingRuleRead(d, m)
}

func resourceRoutingRuleDelete(d *schema.ResourceData, m interface{}) error {
	nexusClient := m.(nexus.Client)

	if err := nexusClient.RoutingRuleDelete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceRoutingRuleExists(d *schema.ResourceData, m interface{}) (bool, error) {
	nexusClient := m.(nexus.Client)

	rule, err := nexusClient.RoutingRuleRead(d.Id())
	return rule != nil, err
}
