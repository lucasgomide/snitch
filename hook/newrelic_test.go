package hook

import (
	"testing"

	"github.com/lucasgomide/snitch/types"
	"gopkg.in/jarcoal/httpmock.v1"
)

var newrelicHost = "https://api.newrelic.com"

func TestNewRelicDeploySuccessful(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var deploys []types.Deploy
	deploys = append(deploys, types.Deploy{"app-sample", "12345678909", "sha1", "user@g.com", "v15"})

	s := NewRelic{newrelicHost, "app-id-here", "api-key-here", "revision-here"}

	httpmock.RegisterResponder("POST", s.Host+"/v2/applications/"+s.ApplicationId+"/deployments.json",
		httpmock.NewStringResponder(201, `ok`))

	if err := s.CallHook(deploys); err != nil {
		t.Error(err)
	}
}

func TestNewRelicReturnsErrorWhenCreateDeployFails(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var deploys []types.Deploy
	deploys = append(deploys, types.Deploy{"app-sample", "12345678909", "sha1", "user@g.com", "v15"})

	s := NewRelic{newrelicHost, "app-id-here", "api-key-here", "revision-here"}

	httpmock.RegisterResponder("POST", s.Host+"/v2/applications/"+s.ApplicationId+"/deployments.json",
		httpmock.NewStringResponder(502, `error`))

	if err := s.CallHook(deploys); err == nil {
		t.Error("Expected returns error, got no error")
	} else if err.Error() != "NewRelic::CreateDeploy - response status code isn't 201" {
		t.Error(err)
	}
}

func TestNewRelicValidateFields(t *testing.T) {
	s := NewRelic{}

	if err = s.ValidatesFields(); err == nil {
		t.Error("Expected returns error, got nil error")
	} else if err.Error() != "Field host into NewRelic hook is required" {
		t.Error("Expected error Field host into NewRelic hook is required, got", err.Error())
	}
	s.Host = "http://abc"

	if err = s.ValidatesFields(); err == nil {
		t.Error("Expected returns error, got nil error")
	} else if err.Error() != "Field application_id into NewRelic hook is required" {
		t.Error("Expected error Field application_id into NewRelic hook is required, got", err.Error())
	}
	s.ApplicationId = "app-id-here"

	if err = s.ValidatesFields(); err == nil {
		t.Error("Expected returns error, got nil error")
	} else if err.Error() != "Field api_key into NewRelic hook is required" {
		t.Error("Expected error Field api_key into NewRelic hook is required, got", err.Error())
	}
	s.ApiKey = "api-key-here"

	if err = s.ValidatesFields(); err == nil {
		t.Error("Expected returns error, got nil error")
	} else if err.Error() != "Field revision into NewRelic hook is required" {
		t.Error("Expected error Field revision into NewRelic hook is required, got", err.Error())
	}
	s.Revision = "revision-here"

	if err = s.ValidatesFields(); err != nil {
		t.Error("Expected returns no error, got", err.Error())
	}
}
