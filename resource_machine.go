package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/joyent/gosdc/cloudapi"
	"regexp"
	"time"
)

var (
	machineStateRunning = "running"
	machineStateStopped = "stopped"

	machineStateChangeTimeout       = 10 * time.Minute
	machineStateChangeCheckInterval = 10 * time.Second
)

func resourceMachine() *schema.Resource {
	return &schema.Resource{
		Create: wrapCallback(resourceMachineCreate),
		Exists: wrapExistsCallback(resourceMachineExists),
		Read:   wrapCallback(resourceMachineRead),
		Update: wrapCallback(resourceMachineUpdate),
		Delete: wrapCallback(resourceMachineDelete),

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description:  "friendly name",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: resourceMachineValidateName,
			},
			"type": &schema.Schema{
				Description: "machine type (smartmachine or virtualmachine)",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"state": &schema.Schema{
				Description: "current state of the machine",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"dataset": &schema.Schema{
				Description: "dataset URN the machine was provisioned with",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"memory": &schema.Schema{
				Description: "amount of memory the machine has (in Mb)",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"disk": &schema.Schema{
				Description: "amount of disk the machine has (in Gb)",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"ips": &schema.Schema{
				Description: "IP addresses the machine has",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"metadata": &schema.Schema{
				Description: "machine metadata, e.g. authorized-keys",
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true, // TODO: remove when Update is added
			},
			"tags": &schema.Schema{
				Description: "machine tags",
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true, // TODO: remove when Update is added
			},
			"created": &schema.Schema{
				Description: "when the machine was created",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated": &schema.Schema{
				Description: "when the machine was update",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"package": &schema.Schema{
				Description: "name of the pakcage to use on provisioning",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true, // TODO: remove when Update is added
				// TODO: validate that the package is available
			},
			"image": &schema.Schema{
				Description: "image UUID",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true, // TODO: remove when Update is added
				// TODO: validate that the UUID is valid
			},
			"primaryip": &schema.Schema{
				Description: "the primary (public) IP address for the machine",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"networks": &schema.Schema{
				Description: "desired network IDs",
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true, // TODO: remove when Update is added
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				// Default:     []string{"public", "private"},
				// TODO: validate that a valid network is presented
			},
			// TODO: firewall_enabled
		},
	}
}

func resourceMachineCreate(d ResourceData, config *Config) error {
	api, err := config.Cloud()
	if err != nil {
		return err
	}

	var networks []string
	for _, network := range d.Get("networks").([]interface{}) {
		networks = append(networks, network.(string))
	}

	metadata := map[string]string{}
	for k, v := range d.Get("metadata").(map[string]interface{}) {
		metadata[k] = v.(string)
	}

	tags := map[string]string{}
	for k, v := range d.Get("tags").(map[string]interface{}) {
		tags[k] = v.(string)
	}

	machine, err := api.CreateMachine(cloudapi.CreateMachineOpts{
		Name:            d.Get("name").(string),
		Package:         d.Get("package").(string),
		Image:           d.Get("image").(string),
		Networks:        networks,
		Metadata:        metadata,
		Tags:            tags,
		FirewallEnabled: true, // TODO: turn this into another schema field
	})
	if err != nil {
		return err
	}

	err = waitForMachineState(api, machine.Id, machineStateRunning, machineStateChangeTimeout)
	if err != nil {
		return err
	}

	// refresh state after it provisions
	machine, err = api.GetMachine(machine.Id)
	if err != nil {
		return err
	}

	setFromMachine(d, machine)

	return nil
}

func resourceMachineExists(d ResourceData, config *Config) (bool, error) {
	api, err := config.Cloud()
	if err != nil {
		return false, err
	}

	machine, err := api.GetMachine(d.Id())

	return machine != nil && err == nil, err
}

func resourceMachineRead(d ResourceData, config *Config) error {
	api, err := config.Cloud()
	if err != nil {
		return err
	}

	machine, err := api.GetMachine(d.Id())
	if err != nil {
		return err
	}

	setFromMachine(d, machine)

	return nil
}

func resourceMachineUpdate(d ResourceData, config *Config) error {
	api, err := config.Cloud()
	if err != nil {
		return err
	}

	d.Partial(true)

	if d.HasChange("name") {
		if err := api.RenameMachine(d.Id(), d.Get("name").(string)); err != nil {
			return err
		}

		err := waitFor(
			func() (bool, error) {
				machine, err := api.GetMachine(d.Id())
				return machine.Name == d.Get("name").(string), err
			},
			machineStateChangeCheckInterval,
			1*time.Minute,
		)
		if err != nil {
			return err
		}

		d.SetPartial("name")
	}

	d.Partial(false)

	return nil
}

func resourceMachineDelete(d ResourceData, config *Config) error {
	api, err := config.Cloud()
	if err != nil {
		return err
	}

	state, err := readMachineState(api, d.Id())
	if state != machineStateStopped {
		err = api.StopMachine(d.Id())
		if err != nil {
			return err
		}

		waitForMachineState(api, d.Id(), machineStateStopped, machineStateChangeTimeout)
	}

	err = api.DeleteMachine(d.Id())
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func readMachineState(api *cloudapi.Client, id string) (string, error) {
	machine, err := api.GetMachine(id)
	if err != nil {
		return "", err
	}

	return machine.State, nil
}

// waitForMachineState waits for a machine to be in the desired state (waiting
// some seconds between each poll). If it doesn't reach the state within the
// duration specified in `timeout`, it returns ErrMachineStateTimeout.
func waitForMachineState(api *cloudapi.Client, id, state string, timeout time.Duration) error {
	return waitFor(
		func() (bool, error) {
			currentState, err := readMachineState(api, id)
			return currentState == state, err
		},
		machineStateChangeCheckInterval,
		machineStateChangeTimeout,
	)
}

// setFromMachine sets resource data from a machine. This includes the ID.
func setFromMachine(d ResourceData, machine *cloudapi.Machine) {
	d.SetId(machine.Id)
	d.Set("name", machine.Name)
	d.Set("type", machine.Type)
	d.Set("state", machine.State)
	d.Set("dataset", machine.Dataset)
	d.Set("memory", machine.Memory)
	d.Set("disk", machine.Disk)
	d.Set("ips", machine.IPs)
	d.Set("metadata", machine.Metadata)
	d.Set("tags", machine.Tags)
	d.Set("created", machine.Created)
	d.Set("updated", machine.Updated)
	d.Set("package", machine.Package)
	d.Set("image", machine.Image)
	d.Set("primaryip", machine.PrimaryIP)
	d.Set("networks", machine.Networks)
}

func resourceMachineValidateName(value interface{}, name string) (warnings []string, errors []error) {
	warnings = []string{}
	errors = []error{}

	r := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\_\.\-]*$`)
	if !r.Match([]byte(value.(string))) {
		errors = append(errors, fmt.Errorf(`"%s" is not a valid %s`, value.(string), name))
	}

	return warnings, errors
}
