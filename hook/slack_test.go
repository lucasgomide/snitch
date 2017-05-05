package hook

import (
	"github.com/lucasgomide/snitch/types"
	"gopkg.in/jarcoal/httpmock.v1"
	"testing"
)

var webhookUrl = "https://slack.url/123"

func TestWhenNotificatedSuccessful(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", webhookUrl,
		httpmock.NewStringResponder(200, `ok`))

	slack := &Slack{WebhookUrl: webhookUrl}
	var deploy []types.Deploy
	deploy = append(deploy, types.Deploy{App: "app-sample"})

	err := slack.CallHook(deploy)
	if err != nil {
		t.Error(err)
	}
}

func TestWhenResponseStatusCodeIsnt200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", webhookUrl,
		httpmock.NewStringResponder(503, `ok`))

	slack := &Slack{WebhookUrl: webhookUrl}
	var deploy []types.Deploy
	deploy = append(deploy, types.Deploy{App: "app-sample"})

	err := slack.CallHook(deploy)
	if err == nil || err.Error() != "Slack - response status code isn't 200" {
		t.Error("It's Exptected that return error when the response status code isn't 200")
	}
}

func TestReturnsErrorWhenRequestFail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterNoResponder(nil)

	slack := &Slack{WebhookUrl: webhookUrl}
	var deploy []types.Deploy
	deploy = append(deploy, types.Deploy{App: "app-sample"})

	err := slack.CallHook(deploy)
	if err == nil {
		t.Error("The request has been failed but no error was raised")
	}
}
