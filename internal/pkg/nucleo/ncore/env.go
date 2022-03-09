package ncore

import "strings"

type Environment = int

// App Environments
const (
	DevelopmentEnvironment = Environment(iota)
	ProductionEnvironment
	TestingEnvironment
)

func ParseEnvironment(str string) Environment {
	switch strings.ToLower(str) {
	case "p", "prod", "production", "1":
		return ProductionEnvironment
	case "t", "test", "testing", "2":
		return TestingEnvironment
	}
	return DevelopmentEnvironment
}
