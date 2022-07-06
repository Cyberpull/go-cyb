package tests

import (
	"testing"

	"cyberpull.com/go-cyb/validator"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DemoStructData struct {
	Name     string `validator:"required;fieldName:Full Name"`
	Email    string `validator:"required;email;fieldName:Email Address"`
	Accepted bool   `validator:"required"`
	Ignored  bool
}

type ValidatorTestSuite struct {
	suite.Suite

	stringEmptyData   string
	stringDefaultData string
	stringNameData    string
	stringEmailData   string
}

func (s *ValidatorTestSuite) SetupSuite() {
	s.stringEmptyData = ""
	s.stringDefaultData = "default"
	s.stringNameData = "Christian Ezeani"
	s.stringEmailData = "demo@example.com"
}

func (s *ValidatorTestSuite) TestValidateBoolSuccess() {
	info := &validator.Validation{
		Name: "Bool",
		Options: validator.Options{
			Required: validator.OptionValue{
				Value: true,
			},
		},
	}

	err := validator.Validate(true, info)
	require.NoError(s.T(), err)
}

func (s *ValidatorTestSuite) TestValidateBoolError() {
	info := &validator.Validation{
		Name: "Bool",
		Options: validator.Options{
			Required: validator.OptionValue{
				Value: true,
			},
		},
	}

	err := validator.Validate(false, info)
	require.Error(s.T(), err)
}

func (s *ValidatorTestSuite) TestValidateStruct() {
	data := DemoStructData{
		Name:     s.stringNameData,
		Email:    s.stringEmailData,
		Accepted: true,
	}

	err := validator.Validate(data)
	require.NoError(s.T(), err)
}

/********************************************/

func TestValidator(t *testing.T) {
	suite.Run(t, new(ValidatorTestSuite))
}
