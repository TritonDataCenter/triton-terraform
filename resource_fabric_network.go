package main

import (
	"errors"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/joyent/gosdc/cloudapi"
)

func resourceFabricNetwork() *schema.Resource {
	return &schema.Resource{
		Create: wrapCallback(resourceFabricNetworkCreate),
		Exists: wrapExistsCallback(resourceFabricNetworkExists),
		Read:   wrapCallback(resourceFabricNetworkRead),
		Delete: wrapCallback(resourceFabricNetworkDelete),

		Schema: map[string]*schema.Schema{
			"vlan_id": {
				Description: "id of VLAN to create network on",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},

			"name": {
				Description: "name of network",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"public": {
				Description: "flag if network is not a RFC1918 private network",
				Type:        schema.TypeBool,
				Computed:    true,
			},

			"fabric": {
				Description: "flag if network is created on a fabric",
				Type:        schema.TypeBool,
				Computed:    true,
			},

			"description": {
				Description: "description of network",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
			},

			"subnet": {
				Description: "CIDR formatted string describing network",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"start_ip": {
				Description: "first assignable IP on network",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"end_ip": {
				Description: "last assignable IP on network",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"gateway": {
				Description: "address of gateway on network, nat zone is created here",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},

			"resolvers": {
				Description: "list of resolvers to use on network",
				Type:        schema.TypeList,
				// This says it is optional, but the cloud api requires it be present.
				// I'm unsure if this is passing an empty array or object because
				// the error from the api says object value found but array needed.
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
			},

			"routes": {
				Description: "map of static routes for hosts on network",
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
			},

			"internet_nat": {
				Description: "whether to create a NAT to the internet at the gateway IP",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
			},
		},
	}
}

func resourceFabricNetworkCreate(d ResourceData, config *Config) error {
	cloud, err := config.Cloud()
	if err != nil {
		return err
	}

	valid, err := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9_\\\\./-]{1,255}$", d.Get("name").(string))
	if !valid || err != nil {
		return errors.New("\"name\" must be at most 255 characters and contain only letters, numbers, _, \\, /, -, and .")
	}

	var resolvers []string
	for _, resolver := range d.Get("resolvers").([]interface{}) {
		resolvers = append(resolvers, resolver.(string))
	}

	routes := map[string]string{}
	for k, v := range d.Get("routes").(map[string]interface{}) {
		routes[k] = v.(string)
	}

	network, err := cloud.CreateFabricNetwork(
		int16(d.Get("vlan_id").(int)),
		cloudapi.CreateFabricNetworkOpts{
			Name:             d.Get("name").(string),
			Description:      d.Get("description").(string),
			Subnet:           d.Get("subnet").(string),
			ProvisionStartIp: d.Get("start_ip").(string),
			ProvisionEndIp:   d.Get("end_ip").(string),
			Gateway:          d.Get("gateway").(string),
			Resolvers:        resolvers,
			Routes:           routes,
			InternetNAT:      d.Get("internet_nat").(bool),
		})
	if err != nil {
		return err
	}

	d.SetId(network.Id)

	err = resourceFabricNetworkRead(d, config)
	if err != nil {
		return err
	}

	return nil
}

func resourceFabricNetworkExists(d ResourceData, config *Config) (bool, error) {
	cloud, err := config.Cloud()
	if err != nil {
		return false, err
	}

	network, err := cloud.GetFabricNetwork(int16(d.Get("vlan_id").(int)), d.Id())

	return network != nil && err == nil, err
}

func resourceFabricNetworkRead(d ResourceData, config *Config) error {
	cloud, err := config.Cloud()
	if err != nil {
		return err
	}

	network, err := cloud.GetFabricNetwork(int16(d.Get("vlan_id").(int)), d.Id())
	if err != nil {
		return err
	}

	d.SetId(network.Id)
	d.Set("vlan_id", network.VLANId)
	d.Set("name", network.Name)
	d.Set("public", network.Public)
	d.Set("fabric", network.Fabric)
	d.Set("description", network.Description)
	d.Set("subnet", network.Subnet)
	d.Set("start_ip", network.ProvisionStartIp)
	d.Set("end_ip", network.ProvisionEndIp)
	d.Set("gateway", network.Gateway)
	d.Set("resolvers", network.Resolvers)
	d.Set("routes", network.Routes)
	d.Set("internet_nat", network.InternetNAT)

	return nil
}

func resourceFabricNetworkDelete(d ResourceData, config *Config) error {
	cloud, err := config.Cloud()
	if err != nil {
		return err
	}

	err = cloud.DeleteFabricNetwork(int16(d.Get("vlan_id").(int)), d.Id())
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
