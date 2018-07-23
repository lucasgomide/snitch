package hook

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/lucasgomide/snitch/config"
	"github.com/lucasgomide/snitch/types"
	"gopkg.in/jarcoal/httpmock.v1"
)

var (
	h              HookFake
	tsuruFake      TsuruFake
	err            error
	configFilePath = "../testdata/config.yaml"
	d              = []types.Deploy{{"app", "1234125125", "sha1", "user@42.com", "v15"}}
	conf           = map[interface{}]interface{}{"field_sample": "key_value"}
)

type HookFake struct {
	Err         error
	FieldSample string
}

type TsuruFake struct {
	Err error
}

func (h HookFake) CallHook(deploy []types.Deploy) error {
	if h.Err != nil {
		return h.Err
	}
	return nil
}

func (h HookFake) ValidatesFields() error {
	return nil
}

func (t TsuruFake) FindLastDeploy(deploy *[]types.Deploy) error {
	if t.Err != nil {
		return t.Err
	}
	*deploy = append(*deploy, types.Deploy{"app", "1234125125", "sha1", "user@42.com", "v15"})
	return nil
}

func TestHookExecutedSuccessfully(t *testing.T) {
	if err = executeHook(&h, d, conf); err != nil {
		t.Error(err)
	}

	if h.FieldSample == "" {
		t.Error("Expected: FieldSample is not empty, got empty")
	}
}

func TestReturnsErrorWhenCallHookFail(t *testing.T) {
	expected := "CallHook has failed"
	h.Err = errors.New(expected)
	err = executeHook(&h, d, conf)

	if err == nil {
		t.Error("Expected error, got nil")
	} else if err.Error() != expected {
		t.Error("Expected error: " + expected + ", got " + err.Error())
	}
	h.Err = nil
}

func TestSetFieldsWithEnvValues(t *testing.T) {
	os.Setenv("NEW_ENV", "gotham")
	if err = executeHook(&h, d, map[interface{}]interface{}{"field_sample": "$NEW_ENV"}); err != nil {
		t.Error(err)
	}

	if h.FieldSample != "gotham" {
		t.Error("Expected: FieldSample equal to gotham, got", h.FieldSample)
	}
}

func TestShouldExecuteHooksFromConfig(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://dummy.sample",
		httpmock.NewStringResponder(200, `ok`))

	httpmock.RegisterResponder("POST", "http://hangouts.chat.sample",
		httpmock.NewStringResponder(200, `ok`))

	httpmock.RegisterResponder("POST", "https://api.rollbar.com/api/1/deploy/",
		httpmock.NewStringResponder(200, `ok`))

	httpmock.RegisterResponder("POST", "http://sentry.com/api/0/projects/the-answer/for-everything/releases/",
		httpmock.NewStringResponder(201, `ok`))
	httpmock.RegisterResponder("POST", "http://sentry.com/api/0/organizations/the-answer/releases/for-everything-v15/deploys/",
		httpmock.NewStringResponder(201, `ok`))

	httpmock.RegisterResponder("POST", "https://api.newrelic.com/v2/applications/01234/deployments.json",
		httpmock.NewStringResponder(201, `ok`))

	err = config.ReadConfigFile(configFilePath)
	if err != nil {
		t.Error(err)
	}

	var buf bytes.Buffer
	log.SetOutput(&buf)

	Execute(h, tsuruFake)

	msg := buf.String()

	if msg != "" {
		t.Error("Expected that msg is not empty, got empty msg")
		t.Error(msg)
	}
}

func TestShouldLogErrorByHooks(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterNoResponder(nil)

	err = config.ReadConfigFile(configFilePath)
	if err != nil {
		t.Error(err)
	}

	var buf bytes.Buffer
	log.SetOutput(&buf)

	Execute(h, tsuruFake)

	msg := buf.String()

	if msg == "" {
		t.Error("Expected msg is empty, got", msg)
	}
}

func TestReturnsErrorWhenFindLastDeployFail(t *testing.T) {
	expected := "FindLastDeploy has failed"
	tsuruFake.Err = errors.New(expected)

	var buf bytes.Buffer
	log.SetOutput(&buf)

	Execute(h, tsuruFake)

	msg := buf.String()

	if msg == "" {
		t.Error("Expected error, got nil")
	} else if !strings.Contains(msg, expected) {
		t.Error("Expected error: " + expected + ", got " + err.Error())
	}
	tsuruFake.Err = nil
}
