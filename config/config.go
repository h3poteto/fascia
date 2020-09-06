package config

import (
	"io/ioutil"
	"os"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/yaml.v2"
)

func Element(elem string) interface{} {
	env := os.Getenv("APPENV")
	file, err := Assets.Open("/settings.yml")
	if err != nil {
		panic(err)
	}
	buf, err := ioutil.ReadAll(file)
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

func getenv(value, key string) string {
	if len(value) == 0 {
		return os.Getenv(key)
	}
	return value
}

type JwtCustomClaims struct {
	CurrentUserID int64 `json:"current_user_id"`
	jwt.StandardClaims
}
