package nval

func InArrayString(value string, list []string) bool {
	for _, listValue := range list {
		if listValue == value {
			return true
		}
	}
	return false
}

func KeyExists(key interface{}, m map[interface{}]interface{}) bool {
	_, ok := m[key]

	return ok
}
