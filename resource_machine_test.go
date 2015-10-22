package main

import (
	"github.com/joyent/gosdc/cloudapi"
	"github.com/joyent/triton-terraform/helpers"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ResourceMachineSuite struct {
	suite.Suite

	server    *helpers.Server
	config    *Config
	api       *cloudapi.Client
	initialID string
	mock      *MockResourceData
}

func (s *ResourceMachineSuite) SetupTest() {
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

	s.initialID = "aaaaaaaa-bbbb-cccc-dddddddddddd"
	s.mock = NewMockResourceData(
		s.initialID,
		map[string]interface{}{
			"name":     "test",
			"package":  "12345678-aaaa-bbbb-cccc-000000000000",       // Micro
			"image":    "12345678-a1a1-b2b2-c3c3-098765432100",       // SmartOS Std
			"networks": []string{"123abc4d-0011-aabb-2233-ccdd4455"}, // Test-Joyent-Public
			"metadata": map[string]string{
				"metadata.key": "value",
			},
			"tags": map[string]string{
				"tag": "???",
			},
		},
	)
}

func (s *ResourceMachineSuite) TeardownTest() {
	s.server.Stop()
}

func (s *ResourceMachineSuite) TestMachineCreateValid() {
	err := resourceMachineCreate(s.mock, s.config)
	s.Assert().Nil(err)

	// assert that the ID is now set
	s.Assert().NotEqual(s.mock.ID, s.initialID)

	// get the resource back and check the fields
	machine, err := s.api.GetMachine(s.mock.ID)
	s.Require().Nil(err)
	s.Require().NotNil(machine)

	s.Assert().Equal(machine.Name, s.mock.Get("name"))
	s.Assert().Equal(machine.Package, s.mock.Get("package"))
	s.Assert().Equal(machine.Image, s.mock.Get("image"))
	// TODO: the following aren't reflected in the localservices API
	// s.Assert().Equal(machine.Networks, s.mock.Get("networks"))
	// s.Assert().Equal(machine.Metadata, s.mock.Get("metadata"))
	// s.Assert().Equal(machine.Tags, s.mock.Get("tags"))
}

func (s *ResourceMachineSuite) TestMachineCreateInvalid() {
	s.mock.Set("package", "blah")
	err := resourceMachineCreate(s.mock, s.config)
	s.Assert().NotNil(err)
}

func (s *ResourceMachineSuite) TestMachineRead() {
	machine, err := s.api.CreateMachine(cloudapi.CreateMachineOpts{
		Name:    s.mock.Get("name").(string),
		Package: s.mock.Get("package").(string),
		Image:   s.mock.Get("image").(string),
	})
	s.Require().Nil(err)

	s.mock.SetId(machine.Id)
	s.mock.Set("name", "")
	s.mock.Set("package", "")
	s.mock.Set("image", "")

	err = resourceMachineRead(s.mock, s.config)
	s.Assert().Nil(err)

	s.Assert().Equal(s.mock.Get("name"), machine.Name)
	s.Assert().Equal(s.mock.Get("package"), machine.Package)
	s.Assert().Equal(s.mock.Get("image"), machine.Image)
}

func (s *ResourceMachineSuite) TestMachineDelete() {
	machine, err := s.api.CreateMachine(cloudapi.CreateMachineOpts{
		Name:    s.mock.Get("name").(string),
		Package: s.mock.Get("package").(string),
		Image:   s.mock.Get("image").(string),
	})
	s.Require().Nil(err)

	setFromMachine(s.mock, machine)

	err = resourceMachineDelete(s.mock, s.config)
	s.Assert().Nil(err)

	machine, err = s.api.GetMachine(machine.Id)
	s.Assert().Nil(machine)
	s.Assert().NotNil(err)
}

func TestResourceMachineSuite(t *testing.T) {
	suite.Run(t, new(ResourceMachineSuite))
}
