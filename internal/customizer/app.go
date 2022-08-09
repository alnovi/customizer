package customizer

import (
	"alnovi/customizer/internal/customizer/server/action/client_actions"
	"alnovi/customizer/internal/customizer/server/action/client_config_action"
	"alnovi/customizer/internal/customizer/server/action/default_config_actions"
	"alnovi/customizer/internal/customizer/server/middleware"
	"alnovi/customizer/internal/customizer/storage/repository"
	"alnovi/customizer/pkg/logger"
	"alnovi/customizer/pkg/storage/drivers"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"alnovi/customizer/internal/customizer/config"
	"alnovi/customizer/internal/customizer/server"
	"alnovi/customizer/internal/customizer/server/action"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type App struct {
	Config     config.Config
	HttpServer *server.HttpServerAdapter
	Router     *mux.Router
	Storage    *drivers.Mongo
	Cache      *drivers.Redis
	Repository *repositoryStore
}

type repositoryStore struct {
	Client  *repository.ClientRepository
	Config  *repository.ConfigRepository
	Default *repository.DefaultConfigRepository
}

func NewApp(cfg config.Config) *App {
	app := App{
		Config:     cfg,
		Repository: new(repositoryStore),
	}

	app.bootstrapStorage()
	app.bootstrapRepository()
	app.bootstrapRouter()
	app.bootstrapHttpServer()
	app.registerAPIv1Actions()
	app.registerWebActions()

	return &app
}

func (a *App) Run() {
	go func() {
		if err := a.HttpServer.ListenAndServe(); err != nil {
			logger.WithError(err).
				WithField("flag", "app.HttpServer.ListenAndServe").
				Fatal("internal.app.App.Run")
		}
	}()

	logger.WithField("message", "Customizer http server has started").
		WithField("host", a.Config.HttpServer.Host).
		WithField("port", a.Config.HttpServer.Port).
		Info("internal.app.App.Run")

	a.handleClose()
}

func (a *App) bootstrapStorage() {
	var err error

	if a.Storage != nil {
		logger.WithError(errors.New("storage driver already bootstrap")).
			Fatal("internal.app.App.bootstrapStorage")
	}

	client, err := drivers.NewMongo(a.Config.MongoConfig)

	if err != nil {
		logger.WithError(err).
			WithField("flag", "drivers.NewMongo").
			Fatal("internal.app.App.bootstrapStorage")
	}

	err = client.Connect()

	if err != nil {
		logger.WithError(err).
			WithField("flag", "client_actions.Connect").
			Fatal("internal.app.App.bootstrapStorage")
	}

	a.Storage = client
}

func (a *App) bootstrapRepository() {
	a.Repository.Client = repository.NewClientRepository(a.Storage)
	a.Repository.Config = repository.NewConfigRepository(a.Storage)
	a.Repository.Default = repository.NewDefaultConfigRepository(a.Storage)
}

func (a *App) bootstrapRouter() {
	a.Router = mux.NewRouter()
	a.Router.Use(middleware.NewCatchMiddleware().Middleware)
}

func (a *App) bootstrapHttpServer() {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders: []string{"*"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowCredentials: true,
		MaxAge:           600,
	})

	a.HttpServer = server.NewHttpServerAdapter(a.Config.HttpServer)
	a.HttpServer.SetHandler(corsHandler.Handler(a.Router))
}

func (a *App) registerAPIv1Actions() {
	api := a.Router.PathPrefix("/api/v1").Subrouter()

	api.Use(middleware.NewPathParametersMiddleware().Middleware)
	api.Use(middleware.NewResponseHeadersMiddleware().Middleware)

	api.Handle("/clients", client_actions.NewListClientAction(a.Repository.Client)).
		Methods(http.MethodGet)
	api.Handle("/clients", client_actions.NewCreateClientAction(a.Repository.Client)).
		Methods(http.MethodPost)
	api.Handle("/clients/{client_id}", client_actions.NewDeleteClientAction(a.Repository.Client)).
		Methods(http.MethodDelete)

	api.Handle("/clients/{client_id}/static", client_config_action.NewStaticListAction(a.Repository.Client, a.Repository.Config)).
		Methods(http.MethodGet)
	api.Handle("/clients/{client_id}/static/{category}", client_config_action.NewGetStaticConfigAction(a.Repository.Client, a.Repository.Config, a.Repository.Default)).
		Methods(http.MethodGet)
	api.Handle("/clients/{client_id}/static/{category}/normalize", client_config_action.NewGetStaticConfigNormalizeAction(a.Repository.Client, a.Repository.Config, a.Repository.Default)).
		Methods(http.MethodGet)
	api.Handle("/clients/{client_id}/static/{category}", client_config_action.NewStoreStaticConfigAction(a.Repository.Client, a.Repository.Config)).
		Methods(http.MethodPost)
	api.Handle("/clients/{client_id}/static/{category}", client_config_action.NewDeleteStaticCategoryConfigAction(a.Repository.Client, a.Repository.Config)).
		Methods(http.MethodDelete)
	api.Handle("/clients/{client_id}/static/{category}/{namespace}", client_config_action.NewDeleteStaticConfigAction(a.Repository.Client, a.Repository.Config)).
		Methods(http.MethodDelete)

	api.Handle("/default-configs/{collection}", default_config_actions.NewListDefaultConfigAction(a.Repository.Default)).
		Methods(http.MethodGet)
	api.Handle("/default-configs/{collection}/{category}", default_config_actions.NewGetDefaultConfigAction(a.Repository.Default)).
		Methods(http.MethodGet)
	api.Handle("/default-configs/{collection}/{category}", default_config_actions.NewStoreDefaultConfigAction(a.Repository.Default)).
		Methods(http.MethodPost)
	api.Handle("/default-configs/{collection}/{category}/{namespace}", default_config_actions.NewDeleteDefaultConfigAction(a.Repository.Default)).
		Methods(http.MethodDelete)
}

func (a *App) registerWebActions() {
	a.Router.PathPrefix(`/`).Handler(action.NewIndexHandler("/app/web/dist"))
}

func (a *App) handleClose() {
	sigintChan := make(chan os.Signal, 1)

	defer close(sigintChan)

	signal.Notify(sigintChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigintChan

	logger.WithField("message", "Got signal").
		WithField("signal", sig).
		Info("internal.app.App.handleClose")

	a.HttpServer.Close()

	logger.WithField("message", "Http server has stopped").
		Info("internal.app.App.handleClose")

	err := a.Storage.Close()

	if err != nil {
		logger.WithError(err).
			WithField("flag", "a.Storage.Disconnect").
			Error("internal.app.App.handleClose")
	} else {
		logger.WithField("message", "Storage has stopped").
			Info("internal.app.App.handleClose")
	}

	err = a.Cache.Disconnect()

	if err != nil {
		logger.WithError(err).
			WithField("flag", "a.Cache.Disconnect").
			Error("internal.app.App.handleClose")
	} else {
		logger.WithField("message", "Cache has stopped").
			Info("internal.app.App.handleClose")
	}
}
