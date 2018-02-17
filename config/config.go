package config

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"os"
	"time"
)

func Element(elem string) interface{} {
	env := os.Getenv("APPENV")
	file, err := Assets.Open("/settings.yml")
	if err != nil {
		panic(err)
	}
	by := new(bytes.Buffer)
	io.Copy(by, file)
	buf := by.Bytes()
	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		panic(err)
	}
	return m[env].(map[interface{}]interface{})[elem]
}

// AWS returns a aws config authorized profile, env, or IAMRole
func AWS() *aws.Config {
	return &aws.Config{
		Credentials: newCredentials(),
		Region:      getRegion(),
	}
}

func newCredentials() *credentials.Credentials {
	return credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.SharedCredentialsProvider{},
			&credentials.EnvProvider{},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(session.New(&aws.Config{
					HTTPClient: &http.Client{Timeout: 3000 * time.Millisecond},
				},
				)),
			},
		})
}

func getRegion() *string {
	return aws.String(os.Getenv("AWS_DEFAULT_REGION"))
}
