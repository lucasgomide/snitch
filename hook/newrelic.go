package hook

import (
	"bytes"
	"errors"
	"net/http"
	"time"

	"github.com/lucasgomide/snitch/types"
)

// NewRelic represents a single struct to notify New Relic
// about a new deploy
type NewRelic struct {
	Host          string
	ApplicationId string
	ApiKey        string
	Revision      string
}

// CallHook will notify New Relic about a new deploy
// It returns any errors encountered
func (s NewRelic) CallHook(deploys []types.Deploy) error {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	return s.createDeploy(httpClient, deploys[0])
}

func (s *NewRelic) createDeploy(httpClient *http.Client, deploy types.Deploy) error {
	data := []byte(`
	{
		"deployment": {
			"revision": "` + s.Revision + `",
			"changelog": "",
			"description": "",
			"user": "` + deploy.User + `"
		}
	}`)

	url := s.Host + "/v2/applications/" + s.ApplicationId + "/deployments.json"
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Add("X-Api-Key", s.ApiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)

	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		return errors.New("NewRelic::CreateDeploy - response status code isn't 201")
	}
	return nil
}

// ValidatesFields checks if there are some field on NewRelic struct invalid
// It returns an error if there are some invalid field
// and if there are no, returns nil
func (s NewRelic) ValidatesFields() error {
	if s.Host == "" {
		return errors.New("Field host into NewRelic hook is required")
	}
	if s.ApplicationId == "" {
		return errors.New("Field application_id into NewRelic hook is required")
	}
	if s.ApiKey == "" {
		return errors.New("Field api_key into NewRelic hook is required")
	}
	if s.Revision == "" {
		return errors.New("Field revision into NewRelic hook is required")
	}
	return nil
}
