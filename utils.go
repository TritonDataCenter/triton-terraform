package main

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"os/user"
	"strings"
	"time"
)

var (
	// ErrNoDefault is returned when coalesceToDefault cannot find a default value
	ErrNoDefault = errors.New("could not find a default value")

	// ErrTimeout is returned when waiting for state change
	ErrTimeout = errors.New("timed out while waiting for resource change")
)

// ExpandUser expands a tilde at the beginning of a path to the current user's
// home directory
func ExpandUser(path string) (string, error) {
	if path[:2] != "~/" {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return path, err
	}

	return strings.Replace(path, "~", usr.HomeDir, 1), nil
}

func coalesceToDefault(defaults ...string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		for _, val := range defaults {
			if val != "" {
				return val, nil
			}
		}

		return "", ErrNoDefault
	}
}

func waitFor(f func() (bool, error), every, timeout time.Duration) error {
	start := time.Now()

	for time.Since(start) <= timeout {
		stop, err := f()
		if err != nil {
			return err
		}

		if stop {
			return nil
		}

		time.Sleep(every)
	}

	return ErrTimeout
}
