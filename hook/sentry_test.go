package hook

import (
	"github.com/lucasgomide/snitch/types"
	"gopkg.in/jarcoal/httpmock.v1"
	"strings"
	"testing"
)

var host = "https://sentry.io"

func TestCreateReleaseAndDeploySuccessful(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var deploys []types.Deploy
	deploys = append(deploys, types.Deploy{"app-sample", "12345678909", "sha1", "user@g.com", "v15"})

	s := Sentry{host, "org-slug", "project-slug", "auth-token", "env", ""}

	httpmock.RegisterResponder("POST", s.Host+"/api/0/projects/"+s.OrganizationSlug+"/"+s.ProjectSlug+"/releases/",
		httpmock.NewStringResponder(201, `ok`))

	releaseVersion := s.ProjectSlug + "-" + deploys[0].Image
	httpmock.RegisterResponder("POST", s.Host+"/api/0/organizations/"+s.OrganizationSlug+"/releases/"+releaseVersion+"/deploys/",
		httpmock.NewStringResponder(201, `ok`))

	if err := s.CallHook(deploys); err != nil {
		t.Error(err)
	}
}

func TestReturnsErrorWhenCreateReleaseFails(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var deploys []types.Deploy
	deploys = append(deploys, types.Deploy{"app-sample", "12345678909", "sha1", "user@g.com", "v15"})

	s := Sentry{host, "org-slug", "project-slug", "auth-token", "env", ""}

	httpmock.RegisterResponder("POST", s.Host+"/api/0/projects/"+s.OrganizationSlug+"/"+s.ProjectSlug+"/releases/",
		httpmock.NewStringResponder(502, `error`))

	if err := s.CallHook(deploys); err.Error() != "Sentry::CreateRelease - response status code isn't 201" {
		t.Error(err)
	}
}

func TestReturnsErrorWhenThereIsNoResponseForCreateRelease(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var deploys []types.Deploy
	deploys = append(deploys, types.Deploy{"app-sample", "12345678909", "sha1", "user@g.com", "v15"})

	s := Sentry{host, "org-slug", "project-slug", "auth-token", "env", ""}

	httpmock.RegisterNoResponder(nil)

	if err := s.CallHook(deploys); err == nil {
		t.Error("Expected returns error, got no error")
	} else if !strings.Contains(err.Error(), "no responder found") {
		t.Error("Expected that the returns contain no responder found, got", err.Error())
	}
}

func TestReturnsErrorWhenCreateDeployFails(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var deploys []types.Deploy
	deploys = append(deploys, types.Deploy{"app-sample", "12345678909", "sha1", "user@g.com", "v15"})

	s := Sentry{host, "org-slug", "project-slug", "auth-token", "env", ""}

	httpmock.RegisterResponder("POST", s.Host+"/api/0/projects/"+s.OrganizationSlug+"/"+s.ProjectSlug+"/releases/",
		httpmock.NewStringResponder(201, `ok`))

	releaseVersion := s.ProjectSlug + "-" + deploys[0].Image
	httpmock.RegisterResponder("POST", s.Host+"/api/0/organizations/"+s.OrganizationSlug+"/releases/"+releaseVersion+"/deploys/",
		httpmock.NewStringResponder(502, `error`))

	if err := s.CallHook(deploys); err == nil {
		t.Error("Expected returns error, got no error")
	} else if err.Error() != "Sentry::CreateDeploy - response status code isn't 201" {
		t.Error(err)
	}
}

func TestReturnsErrorWhenThereIsNoResponseForCreateDeploy(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var deploys []types.Deploy
	deploys = append(deploys, types.Deploy{"app-sample", "12345678909", "sha1", "user@g.com", "v15"})

	s := Sentry{host, "org-slug", "project-slug", "auth-token", "env", ""}

	httpmock.RegisterResponder("POST", s.Host+"/api/0/projects/"+s.OrganizationSlug+"/"+s.ProjectSlug+"/releases/",
		httpmock.NewStringResponder(201, `ok`))

	releaseVersion := deploys[0].Image + "-" + s.OrganizationSlug
	httpmock.RegisterResponder("POST", s.Host+"/api/0/organizations/"+s.OrganizationSlug+"/releases/"+releaseVersion+"/deploys/", nil)

	if err := s.CallHook(deploys); err == nil {
		t.Error("Expected returns error, got no error")
	} else if !strings.Contains(err.Error(), "no responder found") {
		t.Error("Expected that the returns contain: no responder found, got", err.Error())
	}
}

func TestValidateFields(t *testing.T) {
	s := Sentry{}
	if err := s.ValidatesFields(); err == nil {
		t.Error("Expected returns error, got nil error")
	} else if err.Error() != "Field host into Sentry hook is required" {
		t.Error("Expected error Field host into Sentry hook is required, got", err.Error())
	}

	s.Host = "http://abc"
	if err = s.ValidatesFields(); err == nil {
		t.Error("Expected returns error, got nil error")
	} else if err.Error() != "Field organization_slug into Sentry hook is required" {
		t.Error("Expected error Field organization_slug into Sentry hook is required, got", err.Error())
	}

	s.OrganizationSlug = "org-slug"
	if err = s.ValidatesFields(); err == nil {
		t.Error("Expected returns error, got nil error")
	} else if err.Error() != "Field project_slug into Sentry hook is required" {
		t.Error("Expected error Field project_slug into Sentry hook is required, got", err.Error())
	}

	s.ProjectSlug = "project-slug"
	if err = s.ValidatesFields(); err == nil {
		t.Error("Expected returns error, got nil error")
	} else if err.Error() != "Field auth_token into Sentry hook is required" {
		t.Error("Expected error Field auth_token into Sentry hook is required, got", err.Error())
	}

	s.AuthToken = "acb902"
	if err = s.ValidatesFields(); err == nil {
		t.Error("Expected returns error, got nil error")
	} else if err.Error() != "Field env into Sentry hook is required" {
		t.Error("Expected error Field env into Sentry hook is required, got", err.Error())
	}

	s.Env = "sandbox"
	if err = s.ValidatesFields(); err != nil {
		t.Error("Expected returns no error, got", err.Error())
	}
}
