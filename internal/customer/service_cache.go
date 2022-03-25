package customer

import (
	"encoding/json"
	"fmt"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
)

func (s *Service) CacheGet(key string) (string, error) {
	result, err := s.redis.Get(key)
	if err != nil {
		s.log.Error("Error when get Cache: %v.", logOption.Format(key), logOption.Error(err))
		return "", err
	}

	return result, nil
}

func (s *Service) CacheSetThenGet(key string, value string, expire int64) (string, error) {
	result, err := s.redis.SetThenGet(key, value, expire)
	if err != nil {
		s.log.Error("error when set cache %v", logOption.Format(key), logOption.Error(err))
		return "", err
	}

	return result, nil
}

func (s *Service) CacheGetJwt(key string) string {
	fullKey := fmt.Sprintf("%s:%s:%s", constant.Prefix, constant.CacheTokenJWT, key)

	token, err := s.CacheGet(fullKey)
	if err != nil {
		s.log.Error("error when get cache JWT: %v", logOption.Format(fullKey), logOption.Error(err))
		return ""
	}

	return token
}

func (s *Service) CacheGetGoldSavings(cif string) (*dto.GoldSavingVO, error) {
	key := fmt.Sprintf("%s:%s:%s", constant.Prefix, constant.CacheGoldSavings, cif)

	// Get cache gold saving
	data, err := s.CacheGet(key)
	if err != nil {
		s.log.Error("error found when get cache", logOption.Error(err))
		return nil, err
	}

	goldSaving := dto.GoldSavingVO{}

	err = json.Unmarshal([]byte(data), &goldSaving)
	if err != nil {
		s.log.Error("error found when unmarshal data", logOption.Error(err))
		return nil, err
	}

	return &goldSaving, nil
}

func (s *Service) CacheSetGoldSavings(id string, goldSaving *dto.GoldSavingVO) error {
	key := fmt.Sprintf("%s:%s:%s", constant.Prefix, constant.CacheGoldSavings, id)

	value, err := json.Marshal(goldSaving)
	if err != nil {
		return err
	}

	expiry := monthsToUnix(2)

	_, err = s.redis.SetThenGet(key, string(value), expiry)
	if err != nil {
		return err
	}

	return nil
}
