package customer_svc

import (
	"fmt"
	"net/http"
	"path"

	_ "github.com/lib/pq"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/repository"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/service"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nvalidate"
)

var log = nlogger.Get()

type API struct {
	*contract.PdsApp
	dataSources repository.DataSourceMap
}

func NewAPI(core *ncore.Core, config contract.Config) API {
	return API{
		PdsApp: &contract.PdsApp{
			Core:   core,
			Config: config,
			Repositories: contract.RepositoryMap{
				Customer:        new(repository.Customer),
				VerificationOTP: new(repository.VerificationOTP),
				OTP:             new(repository.OTP),
				Credential:      new(repository.Credential),
				AccessSession:   new(repository.AccessSession),
				AuditLogin:      new(repository.AuditLogin),
				Verification:    new(repository.Verification),
				UserExternal:    new(repository.UserExternal),
			},
			Services: contract.ServiceMap{
				Auth:     new(service.Auth),
				Customer: new(service.Customer),
				OTP:      new(service.OTP),
				Cache:    new(service.Cache),
			},
		},
		dataSources: repository.NewDataSourceMap(),
	}
}

func (a *API) Boot() error {
	// Set value default configs
	a.Config.LoadFromEnv()

	// Init data sources
	err := a.dataSources.Init(a.Config.DataSources)
	if err != nil {
		return err
	}

	// Init repositories
	err = ncore.InitStruct(&a.Repositories, a.initRepository)
	if err != nil {
		return err
	}

	// Init services
	err = a.InitService()
	if err != nil {
		return err
	}

	nvalidate.Init()

	return nil
}

func (a *API) InitRouter() http.Handler {
	// Init router
	router := nhttp.NewRouter(nhttp.RouterOptions{
		LogRequest: true,
		Debug:      a.Config.Server.Debug,
		TrustProxy: a.Config.Server.TrustProxy,
	})

	// Init handlers
	handlers := initHandler(a)

	// Enable cors
	if a.Config.CORS.Enabled {
		log.Debug("CORS Enabled")
		router.Use(a.Config.CORS.NewMiddleware())
	}

	// Set-up Routes
	setUpRoute(router, handlers)

	// Set-up Static
	staticPath := path.Join(a.WorkDir, "/web/static")
	staticDir := http.Dir(staticPath)
	staticServer := http.FileServer(staticDir)
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", staticServer))

	return router
}

func (a *API) setUpStaticRoute(r *nhttp.Router) {
	staticServer := http.FileServer(http.Dir(path.Join(a.WorkDir, "/web/static")))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", staticServer))
}

func (a *API) initRepository(name string, i interface{}) error {
	// Check interface
	r, ok := i.(repository.Initializer)
	if !ok {
		return fmt.Errorf("repository '%s' does not implement repository.Initializer interface", name)
	}

	// Init repository
	err := r.Init(a.dataSources, a.Repositories)
	if err != nil {
		return err
	}

	log.Debugf("Repositories.%s has been initialized", name)

	return nil
}
