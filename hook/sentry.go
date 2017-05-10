package hook

import (
	"bytes"
	"errors"
	"github.com/lucasgomide/snitch/types"
	"net/http"
	"time"
)

type Sentry struct {
	Host             string
	OrganizationSlug string
	ProjectSlug      string
	AuthToken        string
	Env              string
	ReleaseVersion   string
}

func (s Sentry) CallHook(deploys []types.Deploy) error {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	if err := s.createRelease(httpClient, deploys[0]); err != nil {
		return err
	}
	if err := s.createDeploy(httpClient); err != nil {
		return err
	}
	return nil
}

func (s *Sentry) createRelease(httpClient *http.Client, deploy types.Deploy) error {
	s.ReleaseVersion = s.ProjectSlug + "-" + deploy.Image
	data := []byte(`{
		"version": "` + s.ReleaseVersion + `",
		"ref":"` + deploy.Commit + `",
		"commits": [{"id":"` + deploy.Commit + `"}]
	}`)

	req, err := http.NewRequest("POST", s.Host+"/api/0/projects/"+s.OrganizationSlug+"/"+s.ProjectSlug+"/releases/", bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+s.AuthToken)
	req.Header.Add("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		return errors.New("Sentry::CreateRelease - response status code isn't 201")
	}
	return nil
}

func (s Sentry) createDeploy(httpClient *http.Client) error {
	data := []byte(`{"environment": "` + s.Env + `"}`)

	req, err := http.NewRequest("POST", s.Host+"/api/0/organizations/"+s.OrganizationSlug+"/releases/"+s.ReleaseVersion+"/deploys/", bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+s.AuthToken)
	req.Header.Add("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		return errors.New("Sentry::CreateDeploy - response status code isn't 201")
	}
	return nil
}

func (s Sentry) ValidatesFields() error {
	if s.Host == "" {
		return errors.New("Field host into Sentry hook is required")
	}
	if s.OrganizationSlug == "" {
		return errors.New("Field organization_slug into Sentry hook is required")
	}
	if s.ProjectSlug == "" {
		return errors.New("Field project_slug into Sentry hook is required")
	}
	if s.AuthToken == "" {
		return errors.New("Field auth_token into Sentry hook is required")
	}
	if s.Env == "" {
		return errors.New("Field env into Sentry hook is required")
	}
	return nil
}
