package main

import (
	"bytes"
	"flag"
	"github.com/lucasgomide/snitch/config"
	"gopkg.in/jarcoal/httpmock.v1"
	"log"
	"os"
	"strings"
	"testing"
)

func setUpSuite() {
	os.Setenv("TSURU_APP_TOKEN", "abc123")
	os.Setenv("TSURU_HOST", "http://0.0.0.0")
	os.Setenv("TSURU_APPNAME", "someapp-name-prd")
	flag.Set("c", "testdata/config_fake_hook.yaml")
	httpmock.Activate()
	httpmock.RegisterResponder("GET", "http://0.0.0.0/deploys?app=&limit=1",
		httpmock.NewStringResponder(200, `[{}]`))
}

func tearDownSuite() {
	defer httpmock.DeactivateAndReset()
	os.Unsetenv("TSURU_TARGET")
	os.Unsetenv("TSURU_TOKEN")
	os.Unsetenv("TSURU_APPNAME")
	flag.Set("app-name-contains", "")
}

func TestShouldReadFileConfig(t *testing.T) {
	setUpSuite()
	defer tearDownSuite()

	if config.Data() != nil {
		t.Error("Expected config.data nil, got ", config.Data())
	}

	main()

	if config.Data() == nil {
		t.Error("Expected config.data isn't nil, got nil")
	}
}

func TestLogErrorWhenCannotReadFileConfig(t *testing.T) {
	setUpSuite()
	defer tearDownSuite()
	flag.Set("c", "path/error.yaml")
	var buf bytes.Buffer
	log.SetOutput(&buf)

	main()

	expected := "open path/error.yaml: no such file or directory"
	msg := buf.String()

	if !strings.Contains(msg, expected) {
		t.Error("Expected config.data isn't nil, got nil")
	}
}

func TestReturnsErrorWhenAppNameDoesNotContainsSomething(t *testing.T) {
	setUpSuite()
	defer tearDownSuite()
	appContains := "$#@"
	flag.Set("app-name-contains", appContains)

	expected := "Tsuru App Name does not contains " + appContains
	var buf bytes.Buffer
	log.SetOutput(&buf)

	main()
	msg := buf.String()

	if !strings.Contains(msg, expected) {
		t.Errorf("%#v, wanted %#v", msg, expected)
	}
}

func TestRunHookExecuteWithAppNameContainsFlagFilled(t *testing.T) {
	setUpSuite()
	httpmock.RegisterResponder("GET", "http://0.0.0.0/deploys?app="+os.Getenv("TSURU_APPNAME")+"&limit=1",
		httpmock.NewStringResponder(200, `[{}]`))
	defer tearDownSuite()
	appContains := "prd"
	flag.Set("app-name-contains", appContains)

	noExpected := "Tsuru App Name does not contains " + appContains
	var buf bytes.Buffer
	log.SetOutput(&buf)

	main()
	msg := buf.String()

	if strings.Contains(msg, noExpected) {
		t.Errorf("Expected no error, got %#v", noExpected)
	}
}
