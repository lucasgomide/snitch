package hook

import (
	"testing"

	"github.com/lucasgomide/snitch/types"
	"gopkg.in/jarcoal/httpmock.v1"
)

var hangoutWebhookUrl = "https://hangouts.chat/123"

func TestHangoutWhenNotificatedSuccessful(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", hangoutWebhookUrl,
		httpmock.NewStringResponder(200, `ok`))

	hangout := &HangoutsChat{WebhookUrl: hangoutWebhookUrl}
	var deploy []types.Deploy
	deploy = append(deploy, types.Deploy{App: "app-sample"})

	err := hangout.CallHook(deploy)
	if err != nil {
		t.Error(err)
	}
}

func TestHangoutWhenResponseStatusCodeIsnt200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", hangoutWebhookUrl,
		httpmock.NewStringResponder(400, `error`))

	hangout := &HangoutsChat{WebhookUrl: hangoutWebhookUrl}
	var deploy []types.Deploy
	deploy = append(deploy, types.Deploy{App: "app-sample"})

	err := hangout.CallHook(deploy)
	expected := "HangoutsChat - response status code is 400"
	if err == nil || err.Error() != expected {
		t.Error("Expected: "+expected+", but got", err.Error())
	}
}

func TestHangoutReturnsErrorWhenRequestFail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterNoResponder(nil)

	hangout := &HangoutsChat{WebhookUrl: hangoutWebhookUrl}
	var deploy []types.Deploy
	deploy = append(deploy, types.Deploy{App: "app-sample"})

	err := hangout.CallHook(deploy)
	if err == nil {
		t.Error("The request has been failed but no error was raised")
	}
}
