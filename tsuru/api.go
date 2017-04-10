package tsuru

import (
	"encoding/json"
	"errors"
	"github.com/lucasgomide/snitch/types"
	"io/ioutil"
	"net/http"
	"time"
)

type TsuruAPI struct {
	AppToken string
	AppName  string
	ApiHost  string
}

func (t TsuruAPI) FindLastDeploy(deploy *[]snitch.Deploy) error {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", t.ApiHost+"/deploys?app="+t.AppName+"&limit=1", nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", t.AppToken)
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("TsuruAPI - response status code isn't 200")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, deploy)
}
