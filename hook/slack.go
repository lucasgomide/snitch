package hook

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/lucasgomide/snitch/types"
)

// Slack represents a single struct to notify Slack
// about a new deploy
type Slack struct {
	WebhookUrl string
}

// CallHook creates a new release and deploy on Slack
// It returns any errors encountered
func (s Slack) CallHook(deploy []types.Deploy) error {
	message := `"The application *` + deploy[0].App + `* has been deployed just now by ` + deploy[0].User + ` at _` + deploy[0].ConvertTimestampToRFC822() + `_"`

	data := []byte(`{"text":` + message + `}`)
	resp, err := http.Post(s.WebhookUrl, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("Slack - response status code isn't 200")
	}
	return nil
}

// ValidatesFields checks if there are some field on Slack struct invalid
// It returns an error if there are some invalid field
// and if there are no, returns nil
func (s Slack) ValidatesFields() error {
	return nil
}
