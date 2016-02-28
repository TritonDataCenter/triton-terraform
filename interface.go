package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// ResourceData is an interface for making more testable callbacks to resources.
// It can be used in combination with wrapCallback wherever you would use a
// typical Terraform callback function. Maintainers, note that you can add any
// methods used by schema.ResourceData onto this. If wrapCallback starts giving
// compile errors, you'll know you've added something incorrectly.
type ResourceData interface {
	Id() string
	SetId(string)
	Get(string) interface{}
	GetOk(string) (interface{}, bool)
	Set(string, interface{}) error
	Partial(bool)
	SetPartial(string)
	HasChange(string) bool
}

func wrapCallback(inner func(*schema.ResourceData, *Config) error) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		return inner(d, meta.(*Config))
	}
}

func wrapExistsCallback(inner func(*schema.ResourceData, *Config) (bool, error)) func(*schema.ResourceData, interface{}) (bool, error) {
	return func(d *schema.ResourceData, meta interface{}) (bool, error) {
		return inner(d, meta.(*Config))
	}
}
