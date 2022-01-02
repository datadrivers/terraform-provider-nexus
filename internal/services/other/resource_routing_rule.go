package other

import (
	"strings"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	nexusSchema "github.com/datadrivers/go-nexus-client/nexus3/schema"
	"github.com/datadrivers/terraform-provider-nexus/internal/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceRoutingRule() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create a Nexus Routing Rule.",

		Create: resourceRoutingRuleCreate,
		Read:   resourceRoutingRuleRead,
		Update: resourceRoutingRuleUpdate,
		Delete: resourceRoutingRuleDelete,
		Exists: resourceRoutingRuleExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func getRoutingRuleFromResourceData(d *schema.ResourceData) nexusSchema.RoutingRule {
	return nexusSchema.RoutingRule{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Mode:        nexusSchema.RoutingRuleMode(d.Get("mode").(string)),
		Matchers:    tools.InterfaceSliceToStringSlice(d.Get("matchers").(*schema.Set).List()),
	}
}

func resourceRoutingRuleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)
	rule := getRoutingRuleFromResourceData(d)

	if err := client.RoutingRule.Create(&rule); err != nil {
		return err
	}

	d.SetId(rule.Name)
	return resourceRoutingRuleRead(d, m)
}

func resourceRoutingRuleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	rule, err := client.RoutingRule.Get(d.Id())
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
	d.Set("matchers", tools.StringSliceToInterfaceSlice(rule.Matchers))

	return nil
}

func resourceRoutingRuleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	rule := getRoutingRuleFromResourceData(d)
	if err := client.RoutingRule.Update(&rule); err != nil {
		return err
	}

	return resourceRoutingRuleRead(d, m)
}

func resourceRoutingRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*nexus.NexusClient)

	if err := client.RoutingRule.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceRoutingRuleExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*nexus.NexusClient)

	rule, err := client.RoutingRule.Get(d.Id())
	return rule != nil, err
}
