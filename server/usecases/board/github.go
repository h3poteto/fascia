package board

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/lib/modules/hub"
	"github.com/h3poteto/fascia/lib/modules/logging"
	rediscli "github.com/h3poteto/fascia/lib/modules/redis"
	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/domains/services"
	"github.com/pkg/errors"
)

func FetchGithub(p *project.Project) (bool, error) {
	return services.FetchGithub(p, InjectProjectRepository(), InjectListRepository(), InjectTaskRepository(), InjectRepoRepository())
}

// InjectRedis returns a redis client instance.
func InjectRedis() *redis.Client {
	return rediscli.SharedInstance().Client
}

// GetAllRepositories gets all repositoreis from github related the oauth token.
func GetAllRepositories(oauthToken string) ([]*github.Repository, error) {
	cli := InjectRedis()
	val, err := cli.LRange(oauthToken, 0, -1).Result()
	logging.SharedInstance().MethodInfo("board", "GetAllRepositories").Debugf("redis values: %+v", val)
	if err == nil && len(val) > 0 {
		var res []*github.Repository
		for _, jsonStr := range val {
			jsonBytes := ([]byte)(jsonStr)
			var r github.Repository
			if err := json.Unmarshal(jsonBytes, &r); err != nil {
				return nil, errors.Wrap(err, "Unmarshal error")
			}
			res = append(res, &r)
		}
		return res, nil
	}
	repositories, err := hub.New(oauthToken).AllRepositories()
	if err != nil {
		return nil, err
	}
	go func() {
		for _, repository := range repositories {
			jsonBytes, err := json.Marshal(repository)
			if err != nil {
				logging.SharedInstance().MethodInfo("board", "GetAllRepositories").Error(err)
				return
			}
			err = cli.RPush(oauthToken, string(jsonBytes)).Err()
			if err != nil {
				logging.SharedInstance().MethodInfo("board", "GetAllRepositories").Error(err)
				return
			}
		}
		// TODO: Extend the expire after you implement refresh function.
		err := cli.Expire(oauthToken, 48*time.Hour).Err()
		if err != nil {
			logging.SharedInstance().MethodInfo("board", "GetAllRepositories").Error(err)
			return
		}
		return
	}()
	return repositories, nil
}
