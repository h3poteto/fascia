package config

import (
	"../modules/logging"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Element(elem string) interface{} {
	env := os.Getenv("GOJIENV")
	root := os.Getenv("GOJIROOT")
	path := filepath.Join(root, "config/settings.yml")
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		logging.SharedInstance().MethodInfo("Config", "Element", true).Panic(err)
	}
	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		logging.SharedInstance().MethodInfo("Config", "Element", true).Panic(err)
	}
	return m[env].(map[interface{}]interface{})[elem]
}
