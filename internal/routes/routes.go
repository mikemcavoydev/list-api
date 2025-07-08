package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/mikemcavoydev/list-api/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)

		r.Get("/lists/{id}", app.Middleware.RequireUser(app.ListHandler.HandleGetListById))
		r.Post("/lists", app.Middleware.RequireUser(app.ListHandler.HandleCreateListById))
		r.Put("/lists/{id}", app.Middleware.RequireUser(app.ListHandler.HandleUpdateListById))
		r.Delete("/lists/{id}", app.Middleware.RequireUser(app.ListHandler.HandleDeleteList))
	})

	r.Get("/health", app.HealthCheck)

	r.Post("/users", app.UserHandler.HandleRegisterUser)

	r.Post("/tokens/authenticate", app.TokenHandler.HandleCreateToken)

	return r
}
