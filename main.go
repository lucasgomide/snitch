package main

import (
	"errors"
	"flag"
	"github.com/lucasgomide/snitch/hook"
	"github.com/lucasgomide/snitch/tsuru"
	"github.com/lucasgomide/snitch/types"
	"log"
	"os"
)

var (
	hookName   = flag.String("h", "", "Service that will provide the webhook. e.g: slack")
	webHookURL = flag.String("url", "", "Webhook URL")
)

func init() {
	flag.Parse()
}

func main() {
	var (
		h      snitch.Hook
		t      snitch.Tsuru
		err    error
		deploy []snitch.Deploy
	)

	if err = validateParams(); err != nil {
		log.Fatal(err)
	}

	switch *hookName {
	case "slack":
		h = &hook.Slack{}
	default:
		log.Fatal("The service " + *hookName + "wasn't implemented yet")
	}

	h.SetWebHookURL(*webHookURL)

	t = tsuru.TsuruAPI{AppToken: os.Getenv("TSURU_APP_TOKEN"), ApiHost: os.Getenv("TSURU_HOST"), AppName: os.Getenv("TSURU_APP_NAME")}

	err = findLastDeploy(t, &deploy)
	if err != nil {
		log.Fatal(err)
	}
	err = callHook(h, deploy)
	if err != nil {
		log.Fatal(err)
	}
}

func findLastDeploy(t snitch.Tsuru, deploy *[]snitch.Deploy) error {
	return t.FindLastDeploy(deploy)
}

func callHook(hook snitch.Hook, deploy []snitch.Deploy) error {
	return hook.CallHook(deploy)
}

func validateParams() error {
	if *hookName == "" {
		return errors.New("The option -h is required ")
	} else if *webHookURL == "" {
		return errors.New("The option -url is required")
	}
	return nil
}
