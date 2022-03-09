package ncore

type Core struct {
	Manifest    Manifest    // Application manifest
	Environment Environment // Process environment
	WorkDir     string      // Working directory
	NodeID      string      // Process Node Identifier
}

func (c *Core) GetEnvironmentString() string {
	switch c.Environment {
	case ProductionEnvironment:
		return "Production"
	case TestingEnvironment:
		return "Testing"
	}
	return "Development"
}
