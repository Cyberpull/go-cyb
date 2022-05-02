package tests

import (
	"testing"

	"cyberpull.com/go-cyb/validator"
)

const (
	stringEmptyData   string = ""
	stringDefaultData string = "default"
	stringNameData    string = "Christian Ezeani"
	stringEmailData   string = "demo@example.com"
)

type DemoStructData struct {
	Name     string `validator:"required;fieldName:Full Name"`
	Email    string `validator:"required;email;fieldName:Email Address"`
	Accepted bool   `validator:"required"`
	Ignored  bool
}

func TestValidator_ValidateBoolSuccess(t *testing.T) {
	info := &validator.Validation{
		Name: "Bool",
		Options: validator.Options{
			Required: validator.OptionValue{
				Value: true,
			},
		},
	}

	if err := validator.Validate(true, info); err != nil {
		t.Fatal(err)
	}
}

func TestValidator_ValidateBoolError(t *testing.T) {
	info := &validator.Validation{
		Name: "Bool",
		Options: validator.Options{
			Required: validator.OptionValue{
				Value: true,
			},
		},
	}

	if err := validator.Validate(false, info); err == nil {
		t.Fatal("Expected an error")
	}
}

func TestValidator_ValidateStruct(t *testing.T) {
	data := DemoStructData{
		Name:     stringNameData,
		Email:    stringEmailData,
		Accepted: true,
	}

	if err := validator.Validate(data); err != nil {
		t.Fatal(err)
	}
}
