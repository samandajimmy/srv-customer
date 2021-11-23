package ncore

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/joho/godotenv"
)

func loadConfig(loadEnvFile bool, dest interface{}, envFile string) error {
	// Validate destination
	validatableConfig, ok := dest.(validation.Validatable)
	if !ok {
		return fmt.Errorf("%s: dest does not implement go-ozzo/validation.Validatable interface", namespace)
	}

	// Load config file from .env
	if loadEnvFile {
		err := godotenv.Load(envFile)
		if err != nil {
			return fmt.Errorf("%s: failed to load configuration"+wrappedErrorFmt, namespace, err)
		}
	}

	// Validate config
	err := validatableConfig.Validate()
	if err != nil {
		return fmt.Errorf("%s: invalid configuration"+wrappedErrorFmt, namespace, err)
	}

	return nil
}
