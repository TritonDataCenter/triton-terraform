package main

import (
	"github.com/joyent/gosdc/cloudapi"
	"github.com/joyent/triton-terraform/helpers"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type ResourceKeySuite struct {
	suite.Suite

	server *helpers.Server
	config *Config
	api    *cloudapi.Client
	mock   *MockResourceData
}

func (s *ResourceKeySuite) SetupTest() {
	var err error
	s.server, err = helpers.NewServer()
	s.Require().Nil(err)

	s.config = &Config{
		Account: helpers.TestAccount,
		Key:     helpers.TestKeyFile,
		KeyID:   helpers.TestKeyID,
		URL:     s.server.URL(),
	}

	s.api, err = s.config.Cloud()
	s.Require().Nil(err)

	s.mock = NewMockResourceData(
		"testkey",
		map[string]interface{}{
			"name": "testkey",
			"key":  helpers.TestPublicKeyData,
		},
	)
}

func (s *ResourceKeySuite) TeardownTest() {
	s.server.Stop()
}

func (s *ResourceKeySuite) TestKeyCreate() {
	err := resourceKeyCreate(s.mock, s.config)
	s.Assert().Nil(err)

	key, err := s.api.GetKey("testkey")
	s.Assert().Nil(err)

	// make sure we created the key OK
	s.Assert().Equal(key.Name, "testkey")
	s.Assert().Equal(key.Key, helpers.TestPublicKeyData)

	// make sure we set the resource ID correctly
	s.Assert().Equal(s.mock.ID, key.Name)
}

func (s *ResourceKeySuite) TestKeyCreateComment() {
	s.mock.Set("name", "")
	err := resourceKeyCreate(s.mock, s.config)
	s.Assert().Nil(err)

	s.Assert().Equal(s.mock.Get("name"), "test@localhost")
}

func (s *ResourceKeySuite) TestKeyCreateNoCommentOrName() {
	s.mock.Set("name", "")
	s.mock.Set("key", strings.Join(strings.SplitN(s.mock.Get("key").(string), " ", 3)[:2], " "))

	err := resourceKeyCreate(s.mock, s.config)
	s.Assert().Equal(err, ErrNoKeyComment)
}

func (s *ResourceKeySuite) TestKeyExists() {
	// it doesn't exist because we haven't created it yet, so let's check that
	exists, err := resourceKeyExists(s.mock, s.config)
	s.Assert().Nil(err)
	s.Assert().False(exists)

	// create the key so we can test the positive case
	_, err = s.api.CreateKey(cloudapi.CreateKeyOpts{
		Name: s.mock.Get("name").(string),
		Key:  s.mock.Get("key").(string),
	})
	s.Assert().Nil(err)

	// now it should exist
	exists, err = resourceKeyExists(s.mock, s.config)
	s.Assert().Nil(err)
	s.Assert().True(exists)
}

func (s *ResourceKeySuite) TestKeyRead() {
	// we're using exists for this resource, so we don't have to test if the
	// resource exists in read. Since that's true, we're just going to go straight
	// to reading an existing key
	key, err := s.api.CreateKey(cloudapi.CreateKeyOpts{
		Name: s.mock.Get("name").(string),
		Key:  s.mock.Get("key").(string),
	})

	err = resourceKeyRead(s.mock, s.config)
	s.Assert().Nil(err)
	s.Assert().Equal(s.mock.Get("name").(string), key.Name)
	s.Assert().Equal(s.mock.Get("key").(string), key.Key)
	s.Assert().Equal(s.mock.ID, key.Name)
}

func (s *ResourceKeySuite) TestKeyDelete() {
	_, err := s.api.CreateKey(cloudapi.CreateKeyOpts{
		Name: s.mock.Get("name").(string),
		Key:  s.mock.Get("key").(string),
	})
	s.Assert().Nil(err)

	err = resourceKeyDelete(s.mock, s.config)
	s.Assert().Nil(err)

	_, err = s.api.GetKey(s.mock.Get("name").(string))
	s.Assert().NotNil(err)
}

func TestKeySuite(t *testing.T) {
	suite.Run(t, new(ResourceKeySuite))
}
