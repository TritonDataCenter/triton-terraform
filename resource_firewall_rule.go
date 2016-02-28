package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/joyent/gosdc/cloudapi"
)

func resourceFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: wrapCallback(resourceFirewallRuleCreate),
		Exists: wrapExistsCallback(resourceFirewallRuleExists),
		Read:   wrapCallback(resourceFirewallRuleRead),
		Update: wrapCallback(resourceFirewallRuleUpdate),
		Delete: wrapCallback(resourceFirewallRuleDelete),

		Schema: map[string]*schema.Schema{
			// required
			"rule": {
				Description: "firewall rule text",
				Type:        schema.TypeString,
				Required:    true,
			},

			// optional
			"enabled": {
				Description: "Indicates if the rule is enabled",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				// TODO: this feels like it should be true by default, but the API docs
				// say that the API has it false by default. Do we want to change this
				// default for this resource?
			},
		},
	}
}

func resourceFirewallRuleCreate(d *schema.ResourceData, config *Config) error {
	cloud, err := config.Cloud()
	if err != nil {
		return err
	}

	rule, err := cloud.CreateFirewallRule(cloudapi.CreateFwRuleOpts{
		Rule:    d.Get("rule").(string),
		Enabled: d.Get("enabled").(bool),
	})
	if err != nil {
		return err
	}

	d.SetId(rule.Id)

	err = resourceFirewallRuleRead(d, config)
	if err != nil {
		return err
	}

	return nil
}

func resourceFirewallRuleExists(d *schema.ResourceData, config *Config) (bool, error) {
	api, err := config.Cloud()
	if err != nil {
		return false, err
	}

	rule, err := api.GetFirewallRule(d.Id())

	return rule != nil && err == nil, err
}

func resourceFirewallRuleRead(d *schema.ResourceData, config *Config) error {
	cloud, err := config.Cloud()
	if err != nil {
		return err
	}

	rule, err := cloud.GetFirewallRule(d.Id())
	if err != nil {
		return err
	}

	d.SetId(rule.Id)
	d.Set("rule", rule.Rule)
	d.Set("enabled", rule.Enabled)

	return nil
}

func resourceFirewallRuleUpdate(d *schema.ResourceData, config *Config) error {
	cloud, err := config.Cloud()
	if err != nil {
		return err
	}

	_, err = cloud.UpdateFirewallRule(
		d.Id(),
		cloudapi.CreateFwRuleOpts{
			Rule:    d.Get("rule").(string),
			Enabled: d.Get("enabled").(bool),
		},
	)
	if err != nil {
		return err
	}

	return resourceFirewallRuleRead(d, config)
}

func resourceFirewallRuleDelete(d *schema.ResourceData, config *Config) error {
	cloud, err := config.Cloud()
	if err != nil {
		return err
	}

	err = cloud.DeleteFirewallRule(d.Id())
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
