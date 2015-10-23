package main

import (
	"github.com/BrianHicks/gosdc/cloudapi"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceKey() *schema.Resource {
	return &schema.Resource{
		Create: wrapCallback(resourceKeyCreate),
		Exists: wrapExistsCallback(resourceKeyExists),
		Read:   wrapCallback(resourceKeyRead),
		Delete: wrapCallback(resourceKeyDelete),

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "name of this key",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"key": &schema.Schema{
				Description: "content of public key from disk",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceKeyCreate(d ResourceData, config *Config) error {
	cloud, err := config.Cloud()
	if err != nil {
		return err
	}

	newKey, err := cloud.CreateKey(cloudapi.CreateKeyOpts{
		Name: d.Get("name").(string),
		Key:  d.Get("key").(string),
	})
	if err != nil {
		return err
	}

	d.SetId(newKey.Name)

	return nil
}

func resourceKeyExists(d ResourceData, config *Config) (bool, error) {
	cloud, err := config.Cloud()
	if err != nil {
		return false, err
	}

	key, err := cloud.GetKey(d.Get("name").(string))
	return key != nil && err == nil, err
}

func resourceKeyRead(d ResourceData, config *Config) error {
	cloud, err := config.Cloud()
	if err != nil {
		return err
	}

	key, err := cloud.GetKey(d.Get("name").(string))

	d.SetId(key.Name)
	d.Set("name", key.Name)
	d.Set("key", key.Key)

	return nil
}

func resourceKeyDelete(d ResourceData, config *Config) error {
	cloud, err := config.Cloud()
	if err != nil {
		return err
	}

	err = cloud.DeleteKey(d.Get("name").(string))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
