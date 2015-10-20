package main

import (
	"github.com/stretchr/testify/assert"
	"os/user"
	"testing"
)

func TestExpandUser(t *testing.T) {
	t.Parallel()

	usr, err := user.Current()
	assert.Nil(t, err)

	// expansion
	expanded, err := ExpandUser("~/test")
	assert.Nil(t, err)
	assert.Equal(t, expanded, usr.HomeDir+"/test")

	// no expansion
	expanded, err = ExpandUser("test")
	assert.Nil(t, err)
	assert.Equal(t, expanded, "test")
}
