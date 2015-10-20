package main

import (
	"github.com/joyent/triton-terraform/helpers"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInitNothing(t *testing.T) {
	c := new(Config)
	assert.Equal(t, c.init(), ErrNoKey)
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
