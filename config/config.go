package config

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

type configuration struct {
	data map[interface{}]interface{}
}

var configs configuration

func ReadConfigFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return readBytes(data)
}

func readBytes(data []byte) error {
	var newConfig map[interface{}]interface{}
	err := yaml.Unmarshal(data, &newConfig)
	if err == nil {
		configs.data = newConfig
	}
	return err
}

func Data() map[interface{}]interface{} {
	return configs.data
}
