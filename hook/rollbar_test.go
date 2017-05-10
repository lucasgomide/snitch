package hook

import (
	"github.com/lucasgomide/snitch/types"
	"gopkg.in/jarcoal/httpmock.v1"
	"strings"
	"testing"
)

func TestCreateDeploySuccessfulOnRollbar(t *testing.T) {
	r := Rollbar{"abc", "development"}

	var (
		deploys       []types.Deploy
		urlRollbarApi = "https://api.rollbar.com/api/1/deploy/"
	)

	deploys = append(deploys, types.Deploy{"", "", "sha1", "user@g.com", ""})

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", urlRollbarApi,
		httpmock.NewStringResponder(200, `ok`))

	if err := r.CallHook(deploys); err != nil {
		t.Error(err)
	}
}

func TestReturnsErrorWhenCreateDeployFailsOnRollbar(t *testing.T) {
	r := Rollbar{"abc", "development"}
	var deploys []types.Deploy
	deploys = append(deploys, types.Deploy{"", "", "sha1", "user@g.com", ""})

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterNoResponder(nil)

	if err := r.CallHook(deploys); err == nil {
		t.Error("Expected returns error, got no error")
	} else if !strings.Contains(err.Error(), "no responder") {
		t.Error("Expected that the returns contain: no responder, got", err.Error())
	}
}

func TestReturnsErrorWhenRequestToCreateDeployIsnt200OnRollbar(t *testing.T) {
	r := Rollbar{"abc123", "development"}
	var (
		deploys       []types.Deploy
		urlRollbarApi = "https://api.rollbar.com/api/1/deploy/"
	)
	deploys = append(deploys, types.Deploy{"", "", "sha1", "user@g.com", ""})

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", urlRollbarApi,
		httpmock.NewStringResponder(501, `error`))
	expected := "Rollbar - response status code isn't 200"
	if err := r.CallHook(deploys); err == nil {
		t.Error("Expected returns error, got no error")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected that the returns contain: %s, got %s", expected, err.Error())
	}
}

func TestValidateFieldsOnRollbar(t *testing.T) {
	r := Rollbar{}
	var err error
	if err = r.ValidatesFields(); err == nil {
		t.Error("Expected returns error, got nil error")
	} else if err.Error() != "Field access_token into Rollbar hook is required" {
		t.Error("Expected error Field access_token into Rollbar hook is required, got", err.Error())
	}

	r.AccessToken = "abc123"
	if err = r.ValidatesFields(); err == nil {
		t.Error("Expected returns error, got nil error")
	} else if err.Error() != "Field env into Rollbar hook is required" {
		t.Error("Expected error Field env into Rollbar hook is required, got", err.Error())
	}

	r.Env = "dev"
	if err = r.ValidatesFields(); err != nil {
		t.Error("Expected returns no error, got", err.Error())
	}
}
