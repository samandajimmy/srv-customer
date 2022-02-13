package customer

import (
	"fmt"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
)

func (s *Service) CacheGet(key string) (string, error) {

	result, err := s.redis.Get(key)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s *Service) CacheSetThenGet(key string, value string, expire int64) (string, error) {

	result, err := s.redis.SetThenGet(key, value, expire)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s *Service) CacheGetJwt(key string) string {
	fullKey := fmt.Sprintf("%s:%s:%s", constant.Prefix, constant.CacheTokenJWT, key)

	token, err := s.CacheGet(fullKey)
	if err != nil {
		return ""
	}

	return token
}
