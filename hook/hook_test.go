package hook

import (
	"bytes"
	"errors"
	"github.com/lucasgomide/snitch/config"
	"github.com/lucasgomide/snitch/types"
	"gopkg.in/jarcoal/httpmock.v1"
	"log"
	"os"
	"strings"
	"testing"
)

var (
	h              HookFake
	tsuruFake      TsuruFake
	err            error
	configFilePath = "../testdata/config.yaml"
	d              = []types.Deploy{types.Deploy{"app", "1234125125", "sha1", "user@42.com"}}
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
	*deploy = append(*deploy, types.Deploy{"app", "1234125125", "sha1", "user@42.com"})
	return nil
}

func TestHookExecutedSuccessfully(t *testing.T) {
	if err = ExecuteHook(&h, d, conf); err != nil {
		t.Error(err)
	}

	if h.FieldSample == "" {
		t.Error("Expected: FieldSample is not empty, got empty")
	}
}

func TestReturnsErrorWhenCallHookFail(t *testing.T) {
	expected := "CallHook has failed"
	h.Err = errors.New(expected)
	err = ExecuteHook(&h, d, conf)

	if err == nil {
		t.Error("Expected error, got nil")
	} else if err.Error() != expected {
		t.Error("Expected error: " + expected + ", got " + err.Error())
	}
	h.Err = nil
}

func TestSetFieldsWithEnvValues(t *testing.T) {
	os.Setenv("NEW_ENV", "gotham")
	if err = ExecuteHook(&h, d, map[interface{}]interface{}{"field_sample": "$NEW_ENV"}); err != nil {
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
