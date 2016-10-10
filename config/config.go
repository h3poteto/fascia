package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

func Element(elem string) interface{} {
	env := os.Getenv("GOJIENV")
	buf, err := Asset("settings.yml")
	if err != nil {
		panic(err)
	}
	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		panic(err)
	}
	return m[env].(map[interface{}]interface{})[elem]
}
