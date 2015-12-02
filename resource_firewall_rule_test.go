package main

import (
	"github.com/joyent/gosdc/cloudapi"
	"github.com/joyent/triton-terraform/helpers"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ResourceFirewallRuleSuite struct {
	suite.Suite

	server    *helpers.Server
	config    *Config
	api       *cloudapi.Client
	initialID string
	mock      *MockResourceData
}

func (s *ResourceFirewallRuleSuite) SetupTest() {
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
			"rule":    "FROM any TO tag www ALLOW tcp PORT 80",
			"enabled": true,
		},
	)
}

func (s *ResourceFirewallRuleSuite) TeardownTest() {
	s.server.Stop()
}

func (s *ResourceFirewallRuleSuite) CreateFirewallRule() *cloudapi.FirewallRule {
	rule, err := s.api.CreateFirewallRule(cloudapi.CreateFwRuleOpts{
		Rule:    s.mock.Get("rule").(string),
		Enabled: s.mock.Get("enabled").(bool),
	})
	s.Require().Nil(err)

	return rule
}

func (s *ResourceFirewallRuleSuite) TestCreateValid() {
	err := resourceFirewallRuleCreate(s.mock, s.config)
	s.Assert().Nil(err)

	// assert that the new ID is now set
	s.Assert().NotEqual(s.mock.ID, s.initialID)

	// get the resource back and check the fields
	rule, err := s.api.GetFirewallRule(s.mock.ID)
	s.Require().Nil(err)
	s.Require().NotNil(rule)

	s.Assert().Equal(rule.Enabled, s.mock.Get("enabled"))
	s.Assert().Equal(rule.Rule, s.mock.Get("rule"))
}

func (s *ResourceFirewallRuleSuite) TestRead() {
	rule := s.CreateFirewallRule()

	s.mock.SetId(rule.Id)
	s.mock.Set("rule", "")
	s.mock.Set("enabled", false)

	err := resourceFirewallRuleRead(s.mock, s.config)
	s.Assert().Nil(err)

	s.Assert().Equal(s.mock.Get("rule"), rule.Rule)
	s.Assert().Equal(s.mock.Get("enabled"), rule.Enabled)
}

func (s *ResourceFirewallRuleSuite) TestUpdate() {
	rule := s.CreateFirewallRule()

	s.mock.SetId(rule.Id)
	newRule := "FROM any TO tag www BLOCK tcp PORT 80"
	s.mock.Change("rule", newRule)
	s.mock.Change("enabled", false)

	err := resourceFirewallRuleUpdate(s.mock, s.config)
	s.Assert().Nil(err)

	rule, err = s.api.GetFirewallRule(rule.Id)
	s.Assert().Nil(err)
	s.Assert().Equal(rule.Rule, newRule)
	s.Assert().False(rule.Enabled)
}

func (s *ResourceFirewallRuleSuite) TestDelete() {
	rule := s.CreateFirewallRule()
	s.mock.SetId(rule.Id)

	err := resourceFirewallRuleDelete(s.mock, s.config)
	s.Assert().Nil(err)

	rule, err = s.api.GetFirewallRule(rule.Id)
	s.Assert().Nil(rule)
	s.Assert().NotNil(err)
}

func TestResourceFirewallRuleSuite(t *testing.T) {
	suite.Run(t, new(ResourceFirewallRuleSuite))
}
