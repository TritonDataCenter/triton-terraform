package main

import (
	"github.com/joyent/gocommon/client"
	"github.com/joyent/gosdc/cloudapi"
	"github.com/joyent/gosign/auth"
	"io/ioutil"
	"os"
)

// Config manages state within the provider
type Config struct {
	Account string `mapstructure:"account"`
	Key     string `mapstructure:"key"`
	KeyID   string `mapstructure:"key_id"`
	URL     string `mapstructure:"url"`

	creds *auth.Credentials
}

// coalesce returns the first non-empty string, or an empty string in the
// terminal case
func (c *Config) coalesce(keys ...string) string {
	for _, val := range keys {
		if val != "" {
			return val
		}
	}

	return ""
}

func (c *Config) init() error {
	c.Account = c.coalesce(c.Account, os.Getenv("SDC_ACCOUNT"))
	c.KeyID = c.coalesce(c.KeyID, os.Getenv("SDC_KEY_ID"))
	c.URL = c.coalesce(c.URL, os.Getenv("SDC_URL"), "https://us-west-1.api.joyentcloud.com")

	key, err := ExpandUser(c.coalesce(c.Key, os.Getenv("SDC_KEY"), "~/.ssh/id_rsa"))
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
