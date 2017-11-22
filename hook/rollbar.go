package hook

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/lucasgomide/snitch/types"
)

// Rollbar represents a single struct to notify Rollbar
// about a new deploy
type Rollbar struct {
	AccessToken string
	Env         string
}

// CallHook creates a new deploy on Rollbar
// It returns any errors encountered
func (r Rollbar) CallHook(deploys []types.Deploy) error {
	data := []byte(
		`{
			"access_token": "` + r.AccessToken + `",
			"environment": "` + r.Env + `",
			"revision": "` + deploys[0].Commit + `",
			"local_username": "` + deploys[0].User + `"
		}
	`)

	resp, err := http.Post("https://api.rollbar.com/api/1/deploy/", "application/json", bytes.NewReader(data))

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("Rollbar - response status code isn't 200")
	}
	return nil
}

// ValidatesFields checks if there are some field on Rollbar struct invalid
// It returns an error if there are some invalid field
// and if there are no, returns nil
func (r Rollbar) ValidatesFields() error {
	if r.AccessToken == "" {
		return errors.New("Field access_token into Rollbar hook is required")
	}

	if r.Env == "" {
		return errors.New("Field env into Rollbar hook is required")
	}

	return nil
}
