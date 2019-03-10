package redis

import (
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

// Redis has client for redis.
type Redis struct {
	Client *redis.Client
}

var sharedInstance = New()

// New returns a redis struct.
func New() *Redis {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	number := os.Getenv("REDIS_DB_NUMBER")
	i, _ := strconv.Atoi(number)
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "", // no password set
		DB:       i,
	})
	return &Redis{
		client,
	}
}

// SharedInstance has redis struct.
func SharedInstance() *Redis {
	return sharedInstance
}
