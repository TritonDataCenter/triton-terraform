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
			"package":  "12345678-aaaa-bbbb-cccc-000000000000",            // Micro
			"image":    "12345678-a1a1-b2b2-c3c3-098765432100",            // SmartOS Std
			"networks": []interface{}{"123abc4d-0011-aabb-2233-ccdd4455"}, // Test-Joyent-Public
			"tags": map[string]interface{}{
				"hello": "world",
			},
		},
	)
}

func (s *ResourceMachineSuite) TeardownTest() {
	s.server.Stop()
}

func (s *ResourceMachineSuite) CreateMachine() *cloudapi.Machine {
	machine, err := s.api.CreateMachine(cloudapi.CreateMachineOpts{
		Name:    s.mock.Get("name").(string),
		Package: s.mock.Get("package").(string),
		Image:   s.mock.Get("image").(string),
	})
	s.Require().Nil(err)

	return machine
}

func (s *ResourceMachineSuite) TestCreateValid() {
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
	// s.Assert().Equal(machine.Metadata, s.mock.Get("metadata"))
	// s.Assert().Equal(machine.Tags, s.mock.Get("tags"))
}

func (s *ResourceMachineSuite) TestCreateInvalid() {
	s.mock.Set("package", "blah")
	err := resourceMachineCreate(s.mock, s.config)
	s.Assert().NotNil(err)
}

func (s *ResourceMachineSuite) TestRead() {
	machine := s.CreateMachine()

	s.mock.SetId(machine.Id)
	s.mock.Set("name", "")
	s.mock.Set("package", "")
	s.mock.Set("image", "")

	err := resourceMachineRead(s.mock, s.config)
	s.Assert().Nil(err)

	s.Assert().Equal(s.mock.Get("name"), machine.Name)
	s.Assert().Equal(s.mock.Get("package"), machine.Package)
	s.Assert().Equal(s.mock.Get("image"), machine.Image)
}

func (s *ResourceMachineSuite) TestUpdateName() {
	machine := s.CreateMachine()

	s.mock.SetId(machine.Id)
	newName := "test-2"
	s.mock.Change("name", newName)

	err := resourceMachineUpdate(s.mock, s.config)
	s.Assert().Nil(err)

	machine, err = s.api.GetMachine(machine.Id)
	s.Assert().Nil(err)
	s.Assert().Equal(machine.Name, newName)
}

func (s *ResourceMachineSuite) TestUpdateTags() {
	machine := s.CreateMachine()

	s.mock.SetId(machine.Id)
	newTags := map[string]interface{}{"hello": "St. Louis"}
	s.mock.Change("tags", newTags)

	// TODO: the ReplaceMachineTags call isn't implemented in the localservices
	// API so this test always fails. It's currently set to succeed if it fails,
	// which is obviously not the best test of the actual functionality.
	err := resourceMachineUpdate(s.mock, s.config)
	s.Assert().NotNil(err)

	machine, err = s.api.GetMachine(machine.Id)
	s.Assert().Nil(err)
	s.Assert().NotEqual(machine.Tags, newTags)
}

func (s *ResourceMachineSuite) TestUpdateTagsEmpty() {
	// TODO: the AddMachineTags call isn't implemented in the localservices
	// API so this test always fails. It's currently set to succeed if it fails,
	// which is obviously not the best test of the actual functionality.
	machine := s.CreateMachine()
	_, err := s.api.AddMachineTags(machine.Id, map[string]string{"test": "value"})
	s.Assert().NotNil(err)

	// s.mock.SetId(machine.Id)
	// newTags := map[string]interface{}{}
	// s.mock.Change("tags", newTags)

	// err = resourceMachineUpdate(s.mock, s.config)
	// s.Assert().Nil(err, err.Error())

	// machine, err = s.api.GetMachine(machine.Id)
	// s.Assert().Nil(err)
	// s.Assert().NotEqual(machine.Tags, newTags)
}

func (s *ResourceMachineSuite) TestUpdatePackage() {
	machine := s.CreateMachine()

	s.mock.SetId(machine.Id)
	newPackage := "11223344-1212-abab-3434-aabbccddeeff" // small
	s.mock.Change("package", newPackage)

	err := resourceMachineUpdate(s.mock, s.config)
	s.Assert().Nil(err)

	machine, err = s.api.GetMachine(machine.Id)
	s.Assert().Nil(err)
	s.Assert().Equal(machine.Package, newPackage)
}

func (s *ResourceMachineSuite) TestUpdateMetadatas() {
	// TODO: the UpdateMachineMetadata call isn't implemented in the localservices
	// API so this test always fails. It's currently set to succeed if it fails,
	// which is obviously not the best test of the actual functionality.

	machine := s.CreateMachine()
	s.mock.SetId(machine.Id)

	for schemaName, apiKey := range resourceMachineMetadataKeys {
		newValue := "something else"
		s.mock.Change(schemaName, newValue)

		err := resourceMachineUpdate(s.mock, s.config)
		s.Assert().NotNil(err, err.Error())

		machine, err = s.api.GetMachine(machine.Id)
		s.Assert().Nil(err)
		s.Assert().NotEqual(machine.Metadata[apiKey], newValue)
	}
}

func (s *ResourceMachineSuite) TestDelete() {
	machine := s.CreateMachine()
	s.mock.SetId(machine.Id)

	err := resourceMachineDelete(s.mock, s.config)
	s.Assert().Nil(err)

	machine, err = s.api.GetMachine(machine.Id)
	s.Assert().Nil(machine)
	s.Assert().NotNil(err)
}

func TestResourceMachineSuite(t *testing.T) {
	suite.Run(t, new(ResourceMachineSuite))
}
