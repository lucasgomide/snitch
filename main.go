package main

import (
	"errors"
	"flag"
	"github.com/lucasgomide/snitch/hook"
	"github.com/lucasgomide/snitch/tsuru"
	"github.com/lucasgomide/snitch/types"
	"log"
	"os"
	"strings"
)

var (
	hookName        = flag.String("hook", "", "Service that will provide the webhook. e.g: slack")
	webHookURL      = flag.String("hook-url", "", "Webhook URL")
	appNameContains = flag.String("app-name-contains", "", "Execute webhook if the tsuru app name contains it")
)

func init() {
	flag.Parse()
}

func execute(h snitch.Hook, t snitch.Tsuru, deploy []snitch.Deploy) error {
	var err error
	if *appNameContains != "" && !strings.Contains(os.Getenv("TSURU_APPNAME"), *appNameContains) {
		return errors.New("Tsuru App Name does not contains " + *appNameContains)
	} else {
		h.SetWebHookURL(*webHookURL)

		err = findLastDeploy(t, &deploy)
		if err != nil {
			return err
		}
		err = callHook(h, deploy)
		if err != nil {
			return err
		}
		return nil
	}
}

func main() {
	var (
		h           snitch.Hook
		t           snitch.Tsuru
		err         error
		deploy      []snitch.Deploy
		hookDefined = false
	)

	if err = validateParams(); err != nil {
		printError(err.Error())
	} else {
		switch *hookName {
		case "slack":
			h = &hook.Slack{}
			hookDefined = true
		default:
			printError("The service " + *hookName + " wasn't implemented yet")
		}
		if hookDefined {
			t = tsuru.TsuruAPI{AppToken: os.Getenv("TSURU_APP_TOKEN"), ApiHost: os.Getenv("TSURU_HOST"), AppName: os.Getenv("TSURU_APP_NAME")}

			if err := execute(h, t, deploy); err != nil {
				printError(err.Error())
			}
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
		return errors.New("The option -hook is required")
	} else if *webHookURL == "" {
		return errors.New("The option -hook-url is required")
	}
	return nil
}

func printError(text string) {
	log.Printf(">> Snitch: %s", text)
}
