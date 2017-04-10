package hook

import (
	"bytes"
	"errors"
	"github.com/lucasgomide/snitch/types"
	"net/http"
)

type Slack struct {
	WebhookURL string
}

func (s *Slack) CallHook(deploy []snitch.Deploy) error {
	message := `"The application ` + deploy[0].App + ` has been deployed just now by ` + deploy[0].User + ` at ` + deploy[0].Timestamp + `"`
	data := []byte(`{"text":` + message + `}`)
	resp, err := http.Post(s.WebhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("Slack - response status code isn't 200")
	}
	return nil
}

func (s *Slack) SetWebHookURL(url string) {
	s.WebhookURL = url
}
