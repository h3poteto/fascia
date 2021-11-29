package config

import (
	_ "embed"
	"os"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/yaml.v2"
)

//go:embed settings.yml
var settings []byte

func Element(elem string) interface{} {
	env := os.Getenv("APPENV")
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal(settings, &m)
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

// JwtCustomClaims jwt claim
type JwtCustomClaims struct {
	CurrentUserID int64 `json:"current_user_id"`
	jwt.StandardClaims
}
