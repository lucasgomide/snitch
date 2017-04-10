package tsuru

import (
	"github.com/lucasgomide/snitch/types"
	"gopkg.in/jarcoal/httpmock.v1"
	"testing"
)

var (
	err      error
	apiHost  = "http://10.0.0.0"
	url      = apiHost + "/deploys?app=" + appName + "&limit=1"
	appName  = "app-name"
	apiToken = "abc12"
)

func TestRetunsLastDeployAsJSON(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(200, `[{"App":"app-test","Timestamp":"2017-04-05T15:21:10.556-03:00","User":"douglas.adams@42.com","Origin":"git"}]`))

	var deploy []snitch.Deploy
	tsuru := TsuruAPI{AppToken: apiToken, AppName: appName, ApiHost: apiHost}

	err = tsuru.FindLastDeploy(&deploy)
	if err != nil {
		t.Error(err)
	}

	if deploy == nil {
		t.Error("The pointer is nil")
	}
}

func TestReturnsErrorWhenRequestFails(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterNoResponder(nil)

	var deploy []snitch.Deploy
	tsuru := TsuruAPI{AppToken: apiToken, AppName: appName, ApiHost: apiHost}

	err = tsuru.FindLastDeploy(&deploy)

	if err == nil {
		t.Error("The request has been failed but no error was raised")
	}

	if deploy != nil {
		t.Error("The pointer should be nil, but isn't")
	}
}

func TestReturnsErrorWhenResponseStatusCodeIsnt200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(503, `ok`))

	var deploy []snitch.Deploy
	tsuru := TsuruAPI{AppToken: apiToken, AppName: appName, ApiHost: apiHost}

	err = tsuru.FindLastDeploy(&deploy)

	if err == nil || err.Error() != "TsuruAPI - response status code isn't 200" {
		t.Error("It's exptected that return error when the response status code isn't 200")
	}

	if deploy != nil {
		t.Error("The pointer should be nil, but isn't")
	}
}
