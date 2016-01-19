package main

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/joyent/gosdc/cloudapi"
)

func resourceFabricVLAN() *schema.Resource {
	return &schema.Resource{
		Create: wrapCallback(resourceFabricVLANCreate),
		Exists: wrapExistsCallback(resourceFabricVLANExists),
		Read:   wrapCallback(resourceFabricVLANRead),
		Update: wrapCallback(resourceFabricVLANUpdate),
		Delete: wrapCallback(resourceFabricVLANDelete),

		Schema: map[string]*schema.Schema{
			"vlan_id": {
				Description:  "VLAN Id between 0-4095",
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: resourceFabricValidateVLAN,
			},

			"name": {
				Description:  "Unique name of VLAN",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: resourceFabricValidateName,
			},

			"description": {
				Description: "Description of VLAN",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
		},
	}
}

func resourceFabricVLANCreate(d ResourceData, config *Config) error {
	cloud, err := config.Cloud()
	if err != nil {
		return err
	}

	_, err = cloud.CreateFabricVLAN(cloudapi.FabricVLAN{
		Id:          int16(d.Get("vlan_id").(int)),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	})
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(d.Get("vlan_id").(int)))

	err = resourceFabricVLANRead(d, config)
	if err != nil {
		return err
	}

	return nil
}

func resourceFabricVLANExists(d ResourceData, config *Config) (bool, error) {
	cloud, err := config.Cloud()
	if err != nil {
		return false, err
	}

	vlanId, err := strconv.ParseInt(d.Id(), 10, 16)
	if err != nil {
		return false, err
	}

	vlan, err := cloud.GetFabricVLAN(int16(vlanId))

	return vlan != nil && err == nil, err
}

func resourceFabricVLANRead(d ResourceData, config *Config) error {
	cloud, err := config.Cloud()
	if err != nil {
		return err
	}

	vlanId, err := strconv.ParseInt(d.Id(), 10, 16)
	if err != nil {
		return err
	}

	vlan, err := cloud.GetFabricVLAN(int16(vlanId))
	if err != nil {
		return err
	}

	d.Set("vlan_id", vlan.Id)
	d.Set("name", vlan.Name)
	d.Set("description", vlan.Description)

	return nil
}

func resourceFabricVLANUpdate(d ResourceData, config *Config) error {
	cloud, err := config.Cloud()
	if err != nil {
		return err
	}

	vlanId, err := strconv.ParseInt(d.Id(), 10, 16)
	if err != nil {
		return err
	}

	_, err = cloud.UpdateFabricVLAN(
		cloudapi.FabricVLAN{
			Id:          int16(vlanId),
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		},
	)
	if err != nil {
		return err
	}

	return resourceFabricVLANRead(d, config)
}

func resourceFabricVLANDelete(d ResourceData, config *Config) error {
	cloud, err := config.Cloud()
	if err != nil {
		return err
	}

	vlanId, err := strconv.ParseInt(d.Id(), 10, 16)
	if err != nil {
		return err
	}

	err = cloud.DeleteFabricVLAN(int16(vlanId))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
