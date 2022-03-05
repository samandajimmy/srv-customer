package nredis

import (
	"errors"
	"fmt"
	"github.com/nbs-go/errx"

	"github.com/gomodule/redigo/redis"
	"github.com/nbs-go/nlogger"
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
		log.Errorf("Connecting to redis client is failed err: %s", err)
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
	if err == redis.ErrNil {
		log.Errorf("Key is empty. key: %s", key)
		return "", nil
	} else if err != nil {
		log.Errorf("Cannot get value from key. err: %s", err)
		return "", nil
	}

	return result, nil

}

func (c *Redis) SetThenGet(key string, value string, expire int64) (string, error) {
	// Set value
	_, err := c.redis.Do("SET", key, value)
	if err != nil {
		log.Errorf("Failed to set value. err: %s", err)
		return "", errx.Trace(err)
	}

	// Set expire
	_, err = c.redis.Do("EXPIRE", key, expire)

	if err != nil {
		log.Errorf("Failed to set expire. err: %s", err)
		return "", errx.Trace(err)
	}

	// Get value
	result, err := redis.String(c.redis.Do("GET", key))

	if errors.Is(err, redis.ErrNil) {
		log.Errorf("Key is empty. key: %s", key)
		return "", nil
	} else if err != nil {
		log.Errorf("Cannot get value from key. err: %s", err)
		return "", errx.Trace(err)
	}

	return result, nil
}
