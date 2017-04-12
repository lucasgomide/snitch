package main

import (
	"errors"
	"flag"
	"github.com/lucasgomide/snitch/hook"
	"github.com/lucasgomide/snitch/tsuru"
	"github.com/lucasgomide/snitch/types"
	"log"
	"os"
	"regexp"
)

var (
	hookName        = flag.String("hook", "", "Service that will provide the webhook. e.g: slack")
	webHookURL      = flag.String("hook-url", "", "Webhook URL")
	appNameContains = flag.String("app-name-contains", "", "Execute webhook if the tsuru app name contains it")
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
	tsuruAppName := os.Getenv("TSURU_APP_NAME")

	if match, _ := regexp.MatchString(*appNameContains, tsuruAppName); !match && *appNameContains != "" {
	} else {

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

		t = tsuru.TsuruAPI{AppToken: os.Getenv("TSURU_APP_TOKEN"), ApiHost: os.Getenv("TSURU_HOST"), AppName: tsuruAppName}

		err = findLastDeploy(t, &deploy)
		if err != nil {
			log.Fatal(err)
		}
		err = callHook(h, deploy)
		if err != nil {
			log.Fatal(err)
		}
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
		return errors.New("The option -hook is required ")
	} else if *webHookURL == "" {
		return errors.New("The option -hook-url is required")
	}
	return nil
}
