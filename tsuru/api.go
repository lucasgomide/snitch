package tsuru

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/lucasgomide/snitch/types"
)

type TsuruAPI struct {
	AppToken string
	AppName  string
	Host     string
}

// FindLastDeploy fetch the last app deploy and writes to deploy
// It returns any errors encountered
func (t TsuruAPI) FindLastDeploy(deploy *[]types.Deploy) error {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", t.Host+"/deploys?app="+t.AppName+"&limit=1", nil)
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
