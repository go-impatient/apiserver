package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
)

// Test file is missing
func TestMissingFile(t *testing.T) {
	filename := "test"
	_, err := Init(filename)

	assert.NotNil(t, err)
}

type ConfigTestSuite struct {
	suite.Suite
	ConfDefault ConfYaml
	Conf        ConfYaml
}

func (suite *ConfigTestSuite) SetupTest() {
	var err error
	suite.ConfDefault, err = Init("")
	if err != nil {
		panic("failed to load default config.yml")
	}
	suite.Conf, err = Init("config.yml")
	if err != nil {
		panic("failed to load config.yml from file")
	}
}

func (suite *ConfigTestSuite) TestValidateConfDefault() {
	// Core
	assert.Equal(suite.T(), "", suite.ConfDefault.Core.Address)
	assert.Equal(suite.T(), "8088", suite.ConfDefault.Core.Port)
	assert.Equal(suite.T(), "debug", suite.Conf.Core.Mode)
	// Log

	// Db
}

func (suite *ConfigTestSuite) TestValidateConf() {
	// Core
	assert.Equal(suite.T(), "8088", suite.Conf.Core.Port)
	assert.Equal(suite.T(), "release", suite.Conf.Core.Mode)
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func TestLoadConfigFromEnv(t *testing.T) {
	os.Setenv("GORUSH_CORE_PORT", "9001")

	Conf, err := Init("config.yml")
	if err != nil {
		panic("failed to load config.yml from file")
	}
	assert.Equal(t, "9001", Conf.Core.Port)
}
