package main

import (
	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/joyent/triton-terraform/helpers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInteralValidate(t *testing.T) {
	t.Parallel()

	p := Provider().(*schema.Provider)

	assert.Nil(t, p.InternalValidate())
}

func TestProviderConfigure(t *testing.T) {
	p := Provider().(*schema.Provider)

	raw, err := config.NewRawConfig(map[string]interface{}{
		"account": helpers.TestAccount,
		"key":     helpers.TestKeyFile,
		"key_id":  helpers.TestKeyID,
		"url":     "https://us-east-1.api.joyentcloud.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	providerConfig := terraform.NewResourceConfig(raw)

	assert.Nil(t, p.Meta())
	assert.Nil(t, p.Configure(providerConfig))

	config, ok := p.Meta().(*Config)
	if assert.True(t, ok) {
		assert.NotNil(t, config)
		assert.Equal(t, config.Account, helpers.TestAccount)
		assert.Equal(t, config.Key, helpers.TestKeyFile)
		assert.Equal(t, config.KeyID, helpers.TestKeyID)
		assert.Equal(t, config.URL, "https://us-east-1.api.joyentcloud.com")
	}
}
