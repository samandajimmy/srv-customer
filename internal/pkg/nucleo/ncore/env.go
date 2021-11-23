package ncore

type Environment = int

// App Environments
const (
	DevelopmentEnvironment = Environment(iota)
	ProductionEnvironment
	TestingEnvironment
)
