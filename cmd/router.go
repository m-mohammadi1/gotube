package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *Application) getRouter() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(commonMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(app.Middleware.EnableCors)

	// use handlers here

	router.Post("/register", app.Handler.UsersHandler.Register)
	router.Post("/login", app.Handler.UsersHandler.Login)
	router.Get("/api/google/authorize/callback", app.Handler.ChannelAuthHandler.HandleCallback)
	//router.Post("/api/google/authorize/callback", app.Handler.ChannelAuthHandler.HandleCallback)

	router.Route("/users", func(r chi.Router) {
		r.Use(app.Middleware.Authorize)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Invalid", 422)
			return
		})
	})

	// channel auth routes
	router.Route("/api/channel/auth", func(r chi.Router) {
		r.Use(app.Middleware.Authorize)

		r.Get("/link", app.Handler.ChannelAuthHandler.GenerateAuthLink)
	})

	return router
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
