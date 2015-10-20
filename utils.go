package main

import (
	"os/user"
	"strings"
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
