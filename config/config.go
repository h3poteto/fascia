package config

import (
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
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

// AWS returns a aws config authorized profile, env, or IAMRole
func AWS(region string) *aws.Config {
	defaultConfig := defaults.Get().Config
	cred := newCredentials(getenv(region, "AWS_DEFAULT_REGION"))
	return defaultConfig.WithCredentials(cred).WithRegion(getenv(region, "AWS_DEFAULT_REGION"))
}

func newCredentials(region string) *credentials.Credentials {
	// temporary config to resolve RemoteCredProvider
	tmpConfig := defaults.Get().Config.WithRegion(region)
	tmpHandlers := defaults.Handlers()

	return credentials.NewChainCredentials(
		[]credentials.Provider{
			// Read profile before environment variables
			&credentials.SharedCredentialsProvider{},
			&credentials.EnvProvider{},
			// for IAM Task Role (ECS) and IAM Role
			defaults.RemoteCredProvider(*tmpConfig, tmpHandlers),
		})
}

func getenv(value, key string) string {
	if len(value) == 0 {
		return os.Getenv(key)
	}
	return value
}
