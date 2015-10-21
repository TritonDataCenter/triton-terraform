package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/mapstructure"
	"os"
)

// Provider returns a Terraform provider for Triton
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: coalesceToDefault(os.Getenv("SDC_ACCOUNT")),
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: coalesceToDefault(
					os.Getenv("SDC_KEY"),
					"~/.ssh/id_rsa",
				),
			},
			"key_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: coalesceToDefault(os.Getenv("SDC_KEY_ID")),
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: coalesceToDefault(
					os.Getenv("SDC_URL"),
					"https://us-west-1.api.joyentcloud.com",
				),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"triton_key": resourceKey(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var config Config
	configRaw := d.Get("").(map[string]interface{})
	if err := mapstructure.Decode(configRaw, &config); err != nil {
		return nil, err
	}
	config.init()

	return &config, nil
}
