package main

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/joyent/gosdc/cloudapi"
	"strings"
)

var (
	// ErrNoKeyComment will be returned when the key name cannot be generated from
	// the key comment and is not otherwise specified.
	ErrNoKeyComment = errors.New("no key comment found to use as a name (and none specified)")
)

func resourceKey() *schema.Resource {
	return &schema.Resource{
		Create: wrapCallback(resourceKeyCreate),
		Exists: wrapExistsCallback(resourceKeyExists),
		Read:   wrapCallback(resourceKeyRead),
		Delete: wrapCallback(resourceKeyDelete),

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "name of this key (will be generated from the key comment, if not set and comment present)",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
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

	if d.Get("name").(string) == "" {
		parts := strings.SplitN(d.Get("key").(string), " ", 3)
		if len(parts) == 3 {
			d.Set("name", parts[2])
		} else {
			return ErrNoKeyComment
		}
	}

	_, err = cloud.CreateKey(cloudapi.CreateKeyOpts{
		Name: d.Get("name").(string),
		Key:  d.Get("key").(string),
	})
	if err != nil {
		return err
	}

	err = resourceKeyRead(d, config)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeyExists(d ResourceData, config *Config) (bool, error) {
	cloud, err := config.Cloud()
	if err != nil {
		return false, err
	}

	keys, err := cloud.ListKeys()
	if err != nil {
		return false, err
	}

	for _, key := range keys {
		if key.Name == d.Id() {
			return true, nil
		}
	}

	return false, nil
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
