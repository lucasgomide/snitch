package config

import (
	"testing"
)

func TestFileDoesntOpen(t *testing.T) {
	err := ReadConfigFile("path/fail")
	if err == nil {
		t.Error("Non error when open a file nonexisting")
	}
}

func TestReadTheFileSuccess(t *testing.T) {
	err := ReadConfigFile("../testdata/config.yaml")
	if err != nil {
		t.Error(err)
	}
}

func TestSetDataUnmarshalIntoConfigs(t *testing.T) {
	err := ReadConfigFile("../testdata/config.yaml")
	if err != nil {
		t.Error(err)
	}

	if Data() == nil {
		t.Error("Data is nil")
	}
}
