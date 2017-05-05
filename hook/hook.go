package hook

import (
	"bytes"
	"github.com/lucasgomide/snitch/config"
	"github.com/lucasgomide/snitch/types"
	"github.com/lucasgomide/snitch/utils"
	"reflect"
	"strings"
)

func Execute(h types.Hook, t types.Tsuru) {
	var err error

	var deploy []types.Deploy

	err = findLastDeploy(t, &deploy)
	if err != nil {
		utils.LogError(err.Error())
	} else {
		for hookName, conf := range config.Data() {
			switch strings.Title(hookName.(string)) {
			case "Slack":
				h = &Slack{}
			default:
				continue
			}
			if err := ExecuteHook(h, deploy, conf); err != nil {
				utils.LogError(err.Error())
			}
		}
	}
}

func findLastDeploy(t types.Tsuru, deploy *[]types.Deploy) error {
	return t.FindLastDeploy(deploy)
}

func callHook(hook types.Hook, deploy []types.Deploy) error {
	return hook.CallHook(deploy)
}

func ExecuteHook(h types.Hook, deploy []types.Deploy, conf interface{}) error {
	s := reflect.ValueOf(h).Elem()

	for k, v := range conf.(map[interface{}]interface{}) {
		var valueBuffer bytes.Buffer
		for _, w := range strings.Split(k.(string), "_") {
			valueBuffer.WriteString(strings.Title(w))
		}
		value := s.FieldByName(valueBuffer.String())
		if value.IsValid() && value.CanAddr() {
			value.SetString(v.(string))
		}
	}

	err := callHook(h, deploy)
	if err != nil {
		return err
	}

	return nil
}
