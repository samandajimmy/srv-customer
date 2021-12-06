package service

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nredis"
)

type Cache struct {
	redis    *nredis.Redis
	response *ncore.ResponseMap
}

func (c *Cache) HasInitialized() bool {
	return true
}

func (c *Cache) Init(app *contract.PdsApp) error {
	redisConfig := app.Config.Redis
	c.redis = nredis.NewNucleoRedis(
		redisConfig.RedisScheme,
		redisConfig.RedisHost,
		redisConfig.RedisPort,
		redisConfig.RedisPass,
	)
	c.response = app.Responses
	return nil
}

func (c *Cache) Get(key string) (string, error) {
	result, err := c.redis.Get(key)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Cache) SetThenGet(key string, value string, expire int64) (string, error) {
	result, err := c.redis.SetThenGet(key, value, expire)
	if err != nil {
		return "", err
	}
	return result, nil
}
