package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type configuration struct {
	data map[interface{}]interface{}
}

var configs configuration

// ReadConfigFile reads the yml config file of given string
// and save configs into struct configuration
// To access the config use Data() function
// It returns any errors encountered
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

// Data reads data's struct field
// It returns the configuration file parsead
func Data() map[interface{}]interface{} {
	return configs.data
}
