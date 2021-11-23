package nhttp

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
)

type CORSConfig struct {
	Enabled        bool
	Origins        []string
	AllowedHeaders []string
	AllowedMethods []string
}

func (c CORSConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Origins,
			validation.When(c.Enabled, validation.Required)),
		validation.Field(&c.AllowedMethods,
			validation.When(c.Enabled, validation.Each(
				validation.In(http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch,
					http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace),
			))),
	)
}

func (c *CORSConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Unmarshal origin as map interface
	tmp := make(map[string]interface{})
	err := unmarshal(&tmp)
	if err != nil {
		c.Enabled = false
		return nil
	}

	// Parse enabled value
	c.Enabled = nval.ParseBooleanFallback(tmp["enabled"], false)

	// Parse origin value if enabled
	if c.Enabled {
		c.Origins = nval.ParseStringArrayFallback(tmp["origin"], []string{"*"})
		c.AllowedHeaders = nval.ParseStringArrayFallback(tmp["allowed_headers"], []string{})
		c.AllowedMethods = nval.ParseStringArrayFallback(tmp["allowed_methods"], []string{http.MethodGet,
			http.MethodHead, http.MethodPost, http.MethodPut, http.MethodOptions})
	}

	return nil
}

func (c *CORSConfig) NewMiddleware() mux.MiddlewareFunc {
	// Init middleware
	optionsAllowedHeaders := handlers.AllowedHeaders(c.AllowedHeaders)
	optionsOrigins := handlers.AllowedOrigins(c.Origins)
	optionsMethods := handlers.AllowedMethods(c.AllowedMethods)
	return handlers.CORS(optionsAllowedHeaders, optionsOrigins, optionsMethods)
}
