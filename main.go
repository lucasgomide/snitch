package main

import (
	"flag"
	"os"
	"strings"

	"github.com/lucasgomide/snitch/config"
	"github.com/lucasgomide/snitch/hook"
	"github.com/lucasgomide/snitch/tsuru"
	"github.com/lucasgomide/snitch/types"
	"github.com/lucasgomide/snitch/utils"
)

var (
	appNameContains = flag.String("app-name-contains", "", "Execute webhook if the tsuru app name contains it")
	configFilePath  = flag.String("c", "", "File path of snitch config")
)

func init() {
	flag.Parse()
}

func main() {
	if *configFilePath == "" {
		utils.LogError("Flag -c is required")
	} else {
		if *appNameContains != "" && !strings.Contains(os.Getenv("TSURU_APPNAME"), *appNameContains) {
			utils.LogError("Tsuru App Name does not contains " + *appNameContains)
		} else {
			run()
		}
	}
}

func run() {
	err := config.ReadConfigFile(*configFilePath)
	if err != nil {
		utils.LogError(err.Error())
	} else {
		var (
			h types.Hook
			t types.Tsuru
		)

		t = tsuru.TsuruAPI{
			AppToken: os.Getenv("TSURU_APP_TOKEN"),
			Host:     os.Getenv("TSURU_HOST"),
			AppName:  os.Getenv("TSURU_APPNAME"),
		}

		hook.Execute(h, t)
	}
}
