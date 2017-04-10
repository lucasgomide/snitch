package hook

import (
	"github.com/lucasgomide/snitch/types"
	"gopkg.in/jarcoal/httpmock.v1"
	"testing"
)

var (
	err        error
	webhookUrl = "https://slack.url/123"
)

func TestWhenNotificatedSuccessful(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", webhookUrl,
		httpmock.NewStringResponder(200, `ok`))

	slack := &Slack{WebhookURL: webhookUrl}
	var deploy []snitch.Deploy
	deploy = append(deploy, snitch.Deploy{App: "app-sample"})

	err = slack.CallHook(deploy)
	if err != nil {
		t.Error(err)
	}
}

func TestWhenResponseStatusCodeIsnt200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", webhookUrl,
		httpmock.NewStringResponder(503, `ok`))

	slack := &Slack{WebhookURL: webhookUrl}
	var deploy []snitch.Deploy
	deploy = append(deploy, snitch.Deploy{App: "app-sample"})

	err = slack.CallHook(deploy)
	if err == nil || err.Error() != "Slack - response status code isn't 200" {
		t.Error("It's Exptected that return error when the response status code isn't 200")
	}
}

func TestReturnsErrorWhenRequestFail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterNoResponder(nil)

	slack := &Slack{WebhookURL: webhookUrl}
	var deploy []snitch.Deploy
	deploy = append(deploy, snitch.Deploy{App: "app-sample"})

	err = slack.CallHook(deploy)
	if err == nil {
		t.Error("The request has been failed but no error was raised")
	}
}

func TestSetWebhookURLSuccessful(t *testing.T) {
	slack := &Slack{}
	slack.SetWebHookURL(webhookUrl)

	if slack.WebhookURL != webhookUrl {
		t.Error("WebhookURL's struct field isn't " + webhookUrl)
	}
}
