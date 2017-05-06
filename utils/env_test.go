package utils

import (
	"os"
	"testing"
)

func TestReturnsEnvFound(t *testing.T) {
	os.Setenv("MY_ENV", "value")

	env := GetEnv("$MY_ENV")

	if env != "value" {
		t.Error("Expetect returns env value, got", env)
	}
	os.Setenv("MY_ENV", "")
}
func TestReturnsInputEnvNotFound(t *testing.T) {
	env := GetEnv("$MY_ENV")

	if env != "$MY_ENV" {
		t.Error("Expetect returns env value, got", env)
	}
}
