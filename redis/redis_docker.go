package sample

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	dockertest "github.com/ory/dockertest"
	"net/url"
)

//CreateReidsClient :return redis-client
func CreateRedisClient() (client *redis.Client, err error) {

	var pool *dockertest.Pool

	pool, err = dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("could not connect redis-docker: %v", err)
	}

	u, err := url.Parse(pool.Client.Endpoint())
	if err != nil {
		return nil, fmt.Errorf("could not parse the endpoint: %v", err)
	}

	resource, err := pool.Run("redis", "latest", []string{})
	if err != nil {
		return nil, fmt.Errorf("could not start resource: %v", err)
	}

	hostName := u.Hostname()
	port := resource.GetPort("6379/tcp")

	client = NewRedisClient(hostName, port, 300, 300, 6500)

	if err = pool.Retry(func() error {
		return pingRedis(client)
	}); err != nil {
		return nil, fmt.Errorf("could not connect redis: %v", err)
	}
	if err = pool.Retry(func() error {
		return checkRedisSetDel(client)
	}); err != nil {
		return nil, fmt.Errorf("could not set or del command redis: %v", err)
	}

	return client, nil
}

func pingRedis(client *redis.Client) error {

	var err error
	_, err = client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func checkRedisSetDel(client *redis.Client) error {

	var key string
	var err error
	key = "connection:test"
	if _, err = client.Set(key, 1, 0).Result(); err != nil {
		return err
	}
	if _, err = client.Del(key).Result(); err != nil {
		return err
	}
	return nil
}

//NewRedis create redis cluster client
func NewRedisClient(hostName string, port string, readTimeout int, writeTimeout int, poolSize int) *redis.Client {

	return redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", hostName, port),
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		PoolSize:     poolSize,
	})
}
