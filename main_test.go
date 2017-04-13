package main

import (
	"bytes"
	"errors"
	"flag"
	"github.com/lucasgomide/snitch/types"
	"gopkg.in/jarcoal/httpmock.v1"
	"log"
	"os"
	"strings"
	"testing"
)

type HookFake struct {
	Err error
}

type TsuruFake struct {
	Err error
}

func (h HookFake) CallHook(deploy []snitch.Deploy) error {
	if h.Err != nil {
		return h.Err
	}
	return nil
}

func (h HookFake) SetWebHookURL(url string) {
}

func (t TsuruFake) FindLastDeploy(deploy *[]snitch.Deploy) error {
	if t.Err != nil {
		return t.Err
	}
	return nil
}

func setUpSuite() {
	os.Setenv("TSURU_APP_TOKEN", "abc123")
	os.Setenv("TSURU_HOST", "http://0.0.0.0")
	os.Setenv("TSURU_APP_NAME", "someapp-name")
}

func tearDownSuite() {
	os.Unsetenv("TSURU_TARGET")
	os.Unsetenv("TSURU_TOKEN")
	flag.Set("hook", "")
	flag.Set("hook-url", "")
	flag.Set("app-name-contains", "")
}

func TestExecutedSuccessfully(t *testing.T) {
	setUpSuite()
	defer tearDownSuite()
	flag.Set("hook", "slack")
	flag.Set("hook-url", "http://123")

	var (
		h      HookFake
		tsuru  TsuruFake
		err    error
		deploy []snitch.Deploy
	)

	if err = execute(h, tsuru, deploy); err != nil {
		t.Error(err)
	}
}

func TestReturnsErrorWhenAppNameDoesNotContainsSomething(t *testing.T) {
	setUpSuite()
	defer tearDownSuite()
	appContains := "$#@"
	flag.Set("hook", "slack")
	flag.Set("hook-url", "http://123")
	flag.Set("app-name-contains", appContains)

	var (
		h      HookFake
		tsuru  TsuruFake
		err    error
		deploy []snitch.Deploy
	)

	wanted := "Tsuru App Name does not match with " + appContains
	err = execute(h, tsuru, deploy)

	if err == nil {
		t.Error("Expected error, got nil")
	} else if err.Error() != wanted {
		t.Error("Expected error: " + wanted + ", wanted " + err.Error())
	}
}

func TestReturnsErrorWhenFindLastDeployFail(t *testing.T) {
	setUpSuite()
	defer tearDownSuite()

	flag.Set("hook-url", "http://123")
	flag.Set("hook", "somehook")

	var (
		h      HookFake
		tsuru  TsuruFake
		err    error
		deploy []snitch.Deploy
	)

	expected := "FindLastDeploy has failed"
	tsuru.Err = errors.New(expected)
	err = execute(h, tsuru, deploy)

	if err == nil {
		t.Error("Expected error, got nil")
	} else if err.Error() != expected {
		t.Error("Expected error: " + expected + ", got " + err.Error())
	}
}

func TestReturnsErrorWhenCallHookFail(t *testing.T) {
	setUpSuite()
	defer tearDownSuite()

	flag.Set("hook-url", "http://123")
	flag.Set("hook", "somehook")

	var (
		h      HookFake
		tsuru  TsuruFake
		err    error
		deploy []snitch.Deploy
	)

	expected := "CallHook has failed"
	h.Err = errors.New(expected)
	err = execute(h, tsuru, deploy)

	if err == nil {
		t.Error("Expected error, got nil")
	} else if err.Error() != expected {
		t.Error("Expected error: " + expected + ", got " + err.Error())
	}
}

func TestPrintErrorWhenHookIsntSupported(t *testing.T) {
	setUpSuite()
	defer tearDownSuite()

	flag.Set("hook", "somehook")
	flag.Set("hook-url", "http://foo.bar")

	var buf bytes.Buffer
	log.SetOutput(&buf)

	main()

	expected := "The service somehook wasn't implemented yet"
	msg := buf.String()
	if !strings.Contains(msg, expected) {
		t.Errorf("%#v, wanted %#v", msg, expected)
	}
}

func TestPrintErrorWhenHookURLFlagIsEmpty(t *testing.T) {
	setUpSuite()
	defer tearDownSuite()

	flag.Set("hook", "somehook")

	var buf bytes.Buffer
	log.SetOutput(&buf)

	main()

	expected := "The option -hook-url is required\n"
	msg := buf.String()
	if !strings.Contains(msg, expected) {
		t.Errorf("%#v, wanted %#v", msg, expected)
	}
}

func TestPrintErrorWhenHookFlagIsEmpty(t *testing.T) {
	setUpSuite()
	defer tearDownSuite()

	flag.Set("hook-url", "http://foo.bar")

	var buf bytes.Buffer
	log.SetOutput(&buf)

	main()

	expected := "The option -hook is required\n"
	msg := buf.String()
	if !strings.Contains(msg, expected) {
		t.Errorf("%#v, wanted %#v", msg, expected)
	}
}

func TestPrintErrorWhenHookRequestFail(t *testing.T) {
	setUpSuite()
	defer tearDownSuite()
	flag.Set("hook-url", "http://foo.bar")
	flag.Set("hook", "slack")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterNoResponder(nil)

	var buf bytes.Buffer
	log.SetOutput(&buf)

	main()

	expected := "no responder found"
	msg := buf.String()
	if !strings.Contains(msg, expected) {
		t.Errorf("%#v, wanted %#v", msg, expected)
	}
}
