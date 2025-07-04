package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/mikemcavoydev/list-api/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	r.Get("/lists/{id}", app.ListHandler.HandleGetListById)
	r.Post("/lists", app.ListHandler.HandleCreateListById)
	r.Put("/lists/{id}", app.ListHandler.HandleUpdateListById)
	r.Delete("/lists/{id}", app.ListHandler.HandleDeleteList)

	r.Post("/users", app.UserHandler.HandleRegisterUser)

	return r
}
