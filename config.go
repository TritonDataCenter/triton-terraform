package main

import (
	"errors"
	"github.com/joyent/gocommon/client"
	"github.com/joyent/gosdc/cloudapi"
	"github.com/joyent/gosign/auth"
	"io/ioutil"
)

var (
	// ErrNoKey is returned by init when .Key is blank
	ErrNoKey = errors.New("key not set")
)

// Config manages state within the provider
type Config struct {
	Account string `mapstructure:"account"`
	Key     string `mapstructure:"key"`
	KeyID   string `mapstructure:"key_id"`
	URL     string `mapstructure:"url"`

	creds *auth.Credentials
}

func (c *Config) init() error {
	if c.Key == "" {
		return ErrNoKey
	}

	key, err := ExpandUser(c.Key)
	if err != nil {
		return err
	}
	c.Key = key

	return nil
}

// Creds returns the credentials needed to connect to Triton
func (c *Config) Creds() (*auth.Credentials, error) {
	if c.creds != nil {
		return c.creds, nil
	}

	keyData, err := ioutil.ReadFile(c.Key)
	if err != nil {
		return nil, err
	}

	userauth, err := auth.NewAuth(c.Account, string(keyData), "rsa-sha256")
	if err != nil {
		return nil, err
	}

	// TODO: this is missing MantaKeyId (string) and MantaEndpoint (string). These
	// will need to be added when Manta support is added to the provider.
	c.creds = &auth.Credentials{
		UserAuthentication: userauth,
		SdcKeyId:           c.KeyID,
		SdcEndpoint:        auth.Endpoint{URL: c.URL},
	}

	return c.creds, nil
}

// Cloud returns a configured cloudapi.Client instance
func (c *Config) Cloud() (*cloudapi.Client, error) {
	creds, err := c.Creds()
	if err != nil {
		return nil, err
	}

	return cloudapi.New(client.NewClient(
		c.URL,
		cloudapi.DefaultAPIVersion,
		creds,
		&cloudapi.Logger,
	)), nil
}
