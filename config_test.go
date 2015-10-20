package main

import (
	"github.com/joyent/triton-terraform/helpers"
	"github.com/stretchr/testify/assert"
	"os"
	"os/user"
	"testing"
)

func TestCoalesce(t *testing.T) {
	t.Parallel()

	config := new(Config)

	// first
	assert.Equal(
		t,
		config.coalesce("test", ""),
		"test",
	)

	// second
	assert.Equal(
		t,
		config.coalesce("", "test"),
		"test",
	)

	// none
	assert.Equal(
		t,
		config.coalesce(""),
		"",
	)
}

func TestExpandPath(t *testing.T) {
	t.Parallel()

	config := new(Config)

	usr, err := user.Current()
	assert.Nil(t, err)

	// expansion
	expanded, err := config.expandPath("~/test")
	assert.Nil(t, err)
	assert.Equal(t, expanded, usr.HomeDir+"/test")

	// no expansion
	expanded, err = config.expandPath("test")
	assert.Nil(t, err)
	assert.Equal(t, expanded, "test")
}

func TestInitNothing(t *testing.T) {
	c := new(Config)
	assert.Nil(t, c.init())
}

func TestInitEnvironment(t *testing.T) {
	unsetAccount := helpers.SetUnset("SDC_ACCOUNT", helpers.TestAccount)
	defer unsetAccount()

	unsetKey := helpers.SetUnset("SDC_KEY", helpers.TestKeyFile)
	defer unsetKey()

	unsetKeyID := helpers.SetUnset("SDC_KEY_ID", helpers.TestKeyID)
	defer unsetKeyID()

	url := "https://us-east-1.api.joyentcloud.com"
	defer helpers.SetUnset("SDC_URL", url)()

	c := new(Config)
	assert.Nil(t, c.init())

	assert.Equal(t, c.Account, helpers.TestAccount)
	assert.Equal(t, c.Key, helpers.TestKeyFile)
	assert.Equal(t, c.KeyID, helpers.TestKeyID)
	assert.Equal(t, c.URL, url)
}

func TestInitHierarchy(t *testing.T) {
	// default with no environment
	c := new(Config)
	assert.Nil(t, c.init())
	assert.Equal(t, c.URL, "https://us-west-1.api.joyentcloud.com")

	// default with environment
	c = new(Config)
	unset := helpers.SetUnset("SDC_URL", "https://us-east-1.api.joyentcloud.com")
	defer unset()
	assert.Nil(t, c.init())
	assert.Equal(t, c.URL, "https://us-east-1.api.joyentcloud.com")

	// explicitly set should not be overridden by either
	c = new(Config)
	c.URL = "test"
	assert.Nil(t, c.init())
	assert.Equal(t, c.URL, "test")
}

func TestCreds(t *testing.T) {
	c := &Config{
		Account: helpers.TestAccount,
		Key:     helpers.TestKeyFile,
		KeyID:   helpers.TestKeyID,
	}
	assert.Nil(t, c.init())

	// get creds
	assert.Nil(t, c.creds)
	creds, err := c.Creds()
	assert.Nil(t, err)
	assert.NotNil(t, creds)

	// try it again, it shouldn't change
	creds2, err := c.Creds()
	assert.Nil(t, err)

	assert.Equal(t, creds, creds2)
	assert.Equal(t, *creds, *creds2)
}

func TestCredsBadFile(t *testing.T) {
	t.Parallel()

	c := &Config{
		Account: helpers.TestAccount,
		Key:     "does.not.exist",
		KeyID:   helpers.TestKeyID,
	}
	assert.Nil(t, c.init())

	// get creds
	assert.Nil(t, c.creds)
	creds, err := c.Creds()
	assert.Nil(t, creds)
	assert.Nil(t, c.creds)
	assert.True(t, os.IsNotExist(err))
}
