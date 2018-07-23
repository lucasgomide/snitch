package hook

import (
	"bytes"
	"reflect"
	"strings"

	"github.com/lucasgomide/snitch/config"
	"github.com/lucasgomide/snitch/types"
	"github.com/lucasgomide/snitch/utils"
)

// Execute runs all hooks that have been configured on
// config file with your options
func Execute(h types.Hook, t types.Tsuru) {
	var err error

	var deploy []types.Deploy

	err = findLastDeploy(t, &deploy)
	if err != nil {
		utils.LogError(err.Error())
	} else {
		for hookName, conf := range config.Data() {
			switch defineHookName(hookName.(string)) {
			case "Slack":
				h = &Slack{}
			case "Sentry":
				h = &Sentry{}
			case "Rollbar":
				h = &Rollbar{}
			case "Newrelic":
				h = &NewRelic{}
			case "Hangouts Chat":
				h = &HangoutsChat{}
			default:
				continue
			}
			if err := executeHook(h, deploy, conf); err != nil {
				utils.LogError(err.Error())
			}
		}
	}
}

func findLastDeploy(t types.Tsuru, deploy *[]types.Deploy) error {
	return t.FindLastDeploy(deploy)
}

func callHook(h types.Hook, deploy []types.Deploy) error {
	err := h.ValidatesFields()
	if err == nil {
		return h.CallHook(deploy)
	}
	return err
}

func executeHook(h types.Hook, deploy []types.Deploy, conf interface{}) error {
	s := reflect.ValueOf(h).Elem()

	for k, v := range conf.(map[interface{}]interface{}) {
		var valueBuffer bytes.Buffer
		for _, w := range strings.Split(k.(string), "_") {
			valueBuffer.WriteString(strings.Title(w))
		}
		value := s.FieldByName(valueBuffer.String())
		if value.IsValid() && value.CanAddr() {
			value.SetString(utils.GetEnv(v.(string)))
		}
	}

	err := callHook(h, deploy)
	if err != nil {
		return err
	}

	return nil
}

func defineHookName(name string) string {
	return strings.Title(
		strings.Join(strings.Split(name, "_"), " "),
	)
}
