package customer

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
