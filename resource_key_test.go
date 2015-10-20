package main

import (
	"github.com/joyent/triton-terraform/helpers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResourceKeyValidateKey(t *testing.T) {
	noWarnings := []string{}
	noErrors := []error{}

	// good key
	warnings, errors := resourceKeyValidateKey(helpers.TestKeyFile, "")
	assert.Equal(t, warnings, noWarnings)
	assert.Equal(t, errors, noErrors)

	// bad key
	warnings, errors = resourceKeyValidateKey("bad.key", "")
	assert.Equal(t, warnings, noWarnings)
	assert.Equal(t, errors, []error{ErrResourceKeyNoPublicKey})
}

func TestResourceKeyCreate(t *testing.T) {
	server, err := helpers.NewServer()
	defer server.Stop()
	assert.Nil(t, err)

	config := &Config{
		Account: helpers.TestAccount,
		Key:     helpers.TestKeyFile,
		KeyID:   helpers.TestKeyID,
		URL:     server.URL(),
	}

	key, err := resourceKeyCreate("test key", helpers.TestPublicKeyFile, config)
	assert.Equal(t, key.Name, "test key")
	assert.Nil(t, err)
}
