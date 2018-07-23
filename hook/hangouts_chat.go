package hook

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/lucasgomide/snitch/types"
)

type HangoutsChat struct {
	WebhookUrl string
}

func (s HangoutsChat) CallHook(deploy []types.Deploy) error {
	message := `"The application *` + deploy[0].App + `* has been deployed just now by ` + deploy[0].User + ` at _` + deploy[0].ConvertTimestampToRFC822() + `_"`

	data := []byte(`{"text":` + message + `}`)
	resp, err := http.Post(s.WebhookUrl, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(`HangoutsChat - response status code is ` + resp.Status)
	}
	return nil
}

func (s HangoutsChat) ValidatesFields() error {
	return nil
}
