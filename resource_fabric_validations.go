package main

import (
	"fmt"
	"regexp"
)

func resourceFabricValidateVLAN(value interface{}, name string) (warnings []string, errors []error) {
	if value.(int) < 0 || value.(int) > 4095 {
		errors = append(errors, fmt.Errorf(`"%s" must be between 0 and 4095`, name))
	}
	return
}

func resourceFabricValidateName(value interface{}, name string) (warnings []string, errors []error) {
	valid, err := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_\\./-]{1,255}$`, value.(string))
	if !valid || err != nil {
		errors = append(errors, fmt.Errorf(`"%s" must be at most 255 characters and contain only letters, numbers, _, \, /, -, and .`, name))
	}
	return
}

func resourceFabricValidateIPv4(value interface{}, name string) (warnings []string, errors []error) {
	valid, err := regexp.MatchString(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`, value.(string))
	if !valid || err != nil {
		errors = append(errors, fmt.Errorf(`"%s" must be an IPv4 address`, name))
	}
	return
}
