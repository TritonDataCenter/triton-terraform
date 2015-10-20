package main

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/joyent/gosdc/cloudapi"
	"io/ioutil"
	"os"
)

var (
	// ErrResourceKeyNoPublicKey is returned when a public key could not be read
	// from disk
	ErrResourceKeyNoPublicKey = errors.New("could not read public key")
)

func resourceKey() *schema.Resource {
	return &schema.Resource{
		Create: func(d *schema.ResourceData, meta interface{}) error {
			key, err := resourceKeyCreate(
				d.Get("name").(string),
				d.Get("key").(string),
				meta.(*Config),
			)
			if err != nil {
				return err
			}

			d.SetId(key.Name)

			return nil
		},
		Exists: resourceKeyExists,
		Read:   resourceKeyRead,
		Delete: resourceKeyDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "name of this key",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"key": &schema.Schema{
				Description:  "public key location on disk",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: resourceKeyValidateKey,
			},
		},
	}
}

func resourceKeyValidateKey(loc interface{}, _ string) (warnings []string, errors []error) {
	warnings = []string{}
	errors = []error{}

	path, err := ExpandUser(loc.(string))
	if err != nil {
		errors = append(errors, err)
		return
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		errors = append(errors, ErrResourceKeyNoPublicKey)
	}

	return
}

func resourceKeyCreate(name, key string, config *Config) (*cloudapi.Key, error) {
	cloud, err := config.Cloud()
	if err != nil {
		return nil, err
	}

	path, err := ExpandUser(key)
	if err != nil {
		return nil, err
	}

	keyData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	newKey, err := cloud.CreateKey(cloudapi.CreateKeyOpts{
		Name: name,
		Key:  string(keyData),
	})
	if err != nil {
		return nil, err
	}

	return newKey, nil
}

func resourceKeyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	config := meta.(*Config)
	cloud, err := config.Cloud()
	if err != nil {
		return false, err
	}

	key, err := cloud.GetKey(d.Get("name").(string))
	return key != nil && err == nil, err
}

func resourceKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cloud, err := config.Cloud()
	if err != nil {
		return err
	}

	key, err := cloud.GetKey(d.Get("name").(string))

	d.Set("name", key.Name)
	d.Set("key", key.Key)

	return nil
}

func resourceKeyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
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
