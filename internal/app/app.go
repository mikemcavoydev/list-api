package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mikemcavoydev/list-api/internal/api"
)

type Application struct {
	Logger      *log.Logger
	ListHandler *api.ListHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	listHandler := api.NewListHandler()

	app := &Application{
		Logger:      logger,
		ListHandler: listHandler,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Status is available\n")
}
