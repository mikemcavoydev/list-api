package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mikemcavoydev/list-api/internal/api"
	"github.com/mikemcavoydev/list-api/internal/middleware"
	"github.com/mikemcavoydev/list-api/internal/store"
	"github.com/mikemcavoydev/list-api/migrations"
)

type Application struct {
	Logger       *log.Logger
	ListHandler  *api.ListHandler
	UserHandler  *api.UserHandler
	TokenHandler *api.TokenHandler
	Middleware   middleware.UserMiddleware
	DB           *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	listStore := store.NewPostgresListStore(pgDB)
	listHandler := api.NewListHandler(listStore, logger)

	userStore := store.NewPostgresUserStore(pgDB)
	userHandler := api.NewUserHandler(userStore, logger)

	tokenStore := store.NewPostgresTokenStore(pgDB)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)

	middlewareHandler := middleware.UserMiddleware{
		UserStore: userStore,
	}

	app := &Application{
		Logger:       logger,
		ListHandler:  listHandler,
		UserHandler:  userHandler,
		TokenHandler: tokenHandler,
		Middleware:   middlewareHandler,
		DB:           pgDB,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Status is available\n")
}
