package customizer

import (
	"fmt"
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
}

func NewApp(cfg config.Config) *App {
	app := App{
		Config: cfg,
	}

	app.bootstrapRouter()
	app.bootstrapHttpServer()

	app.registerWebActions()

	return &app
}

func (a *App) Run() {
	go func() {
		if err := a.HttpServer.ListenAndServe(); err != nil {
			// TODO write log
		}
	}()

	a.handleClose()
}

func (a *App) bootstrapRouter() {
	a.Router = mux.NewRouter()
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

func (a *App) registerWebActions() {
	a.Router.PathPrefix(`/`).Handler(action.NewIndexHandler("./dist"))
}

func (a *App) handleClose() {
	sigintChan := make(chan os.Signal, 1)

	defer close(sigintChan)

	signal.Notify(sigintChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigintChan

	fmt.Printf("Got signal: %v\n", sig)

	a.HttpServer.Close()
}
