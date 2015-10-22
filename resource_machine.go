package main

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/joyent/gosdc/cloudapi"
	"time"
)

var (
	// ErrMachineStateTimeout is returned when changing machine state results in a
	// timeout
	ErrMachineStateTimeout = errors.New("timed out waiting for machine state")

	machineStateRunning = "running"
	machineStateStopped = "stopped"

	machineStateChangeTimeout = 60 * time.Second
)

func resourceMachine() *schema.Resource {
	return &schema.Resource{
		Create: wrapCallback(resourceMachineCreate),
		Exists: wrapExistsCallback(resourceMachineExists),
		Read:   wrapCallback(resourceMachineRead),

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "friendly name",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true, // TODO: remove when Update is added
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
			"networks": &schema.Schema{
				Description: "desired network IDs",
				Type:        schema.TypeList,
				Elem:        schema.TypeString,
				Optional:    true,
				ForceNew:    true, // TODO: remove when Update is added
				// Default:     []string{"public", "private"},
				// TODO: validate that a valid network is presented
			},
			"metadata": &schema.Schema{
				Description: "an arbitrary set of metadata key/value pairs",
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true, // TODO: remove when Update is added
			},
			"tags": &schema.Schema{
				Description: "an arbitrary set of tags",
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true, // TODO: remove when Update is added
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

	machine, err := api.CreateMachine(cloudapi.CreateMachineOpts{
		Name:            d.Get("name").(string),
		Package:         d.Get("package").(string),
		Image:           d.Get("image").(string),
		Networks:        d.Get("networks").([]string),
		Metadata:        d.Get("metadata").(map[string]string),
		Tags:            d.Get("tags").(map[string]string),
		FirewallEnabled: true, // TODO: turn this into another schema field
	})
	if err != nil {
		return err
	}

	err = waitForMachineState(api, machine.Id, machineStateRunning, 60*time.Second)
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

// waitForMachineState waits for a machine to be in the desired state (waiting 5
// seconds between each poll). If it doesn't reach the state within the duration
// specified in `timeout`, it returns ErrMachineStateTimeout
func waitForMachineState(api *cloudapi.Client, id, state string, timeout time.Duration) error {
	start := time.Now()

	for time.Since(start) <= timeout {
		currentState, err := readMachineState(api, id)
		if err != nil {
			return err
		}

		if currentState != state {
			time.Sleep(5 * time.Second)
		} else {
			return nil
		}
	}

	return ErrMachineStateTimeout
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
