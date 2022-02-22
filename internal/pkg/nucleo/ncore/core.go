package ncore

const namespace = "nbs-go/nucleo/ncore"
const wrappedErrorFmt = "\n  > %w"

type Core struct {
	Manifest    Manifest     // Application manifest
	Environment Environment  // Process environment
	WorkDir     string       // Working directory
	NodeID      string       // Process Node Identifier
	Responses   *ResponseMap // Contains list of response codes
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
