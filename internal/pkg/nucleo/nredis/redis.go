package nredis

import (
	"errors"
	"fmt"
	"github.com/nbs-go/errx"
	logOption "github.com/nbs-go/nlogger/v2/option"

	"github.com/gomodule/redigo/redis"
	"github.com/nbs-go/nlogger/v2"
)

var log = nlogger.Get()

type Redis struct {
	rPool *redis.Pool
	redis redis.Conn
}

func NewNucleoRedis(network string, host string, port string, password string) *Redis {
	// Init new pool redis
	pool, err := DialClient(network, host, port, password)
	if err != nil {
		log.Error("Connecting to redis client is failed", logOption.Error(err))
		return nil
	}

	// Get connection from pool
	conn := pool.Get()

	return &Redis{
		rPool: pool,
		redis: conn,
	}
}

func DialClient(network string, host string, port string, password string) (*redis.Pool, error) {
	address := fmt.Sprintf("%s:%s", host, port)

	initPool := redis.Pool{
		MaxIdle:   50,
		MaxActive: 10000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(network, address, redis.DialPassword(password))
			if err != nil {
				return nil, errx.Trace(err)
			}
			return conn, nil
		},
	}

	return &initPool, nil
}

func (c *Redis) Ping() (string, error) {
	result, err := redis.String(c.redis.Do("PING"))
	if err != nil {
		return "", errx.Trace(err)
	}
	return result, nil
}

func (c *Redis) Get(key string) (string, error) {
	result, err := redis.String(c.redis.Do("GET", key))
	if errors.Is(err, redis.ErrNil) {
		log.Errorf("Key is empty. key: %s", key)
		return "", nil
	} else if err != nil {
		log.Error("Cannot get value from key", logOption.Error(err))
		return "", nil
	}

	return result, nil

}

func (c *Redis) SetThenGet(key string, value string, expire int64) (string, error) {
	// Set value
	_, err := c.redis.Do("SET", key, value)
	if err != nil {
		log.Error("Failed to set value", logOption.Error(err))
		return "", errx.Trace(err)
	}

	// Set expire
	_, err = c.redis.Do("EXPIRE", key, expire)

	if err != nil {
		log.Error("Failed to set expire", logOption.Error(err))
		return "", errx.Trace(err)
	}

	// Get value
	result, err := redis.String(c.redis.Do("GET", key))

	if errors.Is(err, redis.ErrNil) {
		log.Errorf("Key is empty. key: %s", key)
		return "", nil
	} else if err != nil {
		log.Error("Cannot get value from key", logOption.Error(err))
		return "", errx.Trace(err)
	}

	return result, nil
}
