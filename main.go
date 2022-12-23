package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/finest08/jwt-auth-demo/handler"
	"github.com/finest08/jwt-auth-demo/store"
)

func main() {
	s := store.Connect()
	r := chi.NewRouter()
	r.Use(
		middleware.Logger,
		middleware.StripSlashes,
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "QUERY"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
	)

	p := &handler.Person{
		Store: s,
	}

	r.Route("/person", func(r chi.Router) {
		r.With().Post("/create", p.Create)
	})

	r.Route("/login", func(r chi.Router) {
		r.With().Post("/", p.Login)
	})

	r.Route("/auth", func(r chi.Router) {
		r.With().Get("/people", p.Query)
		r.With().Get("/user", p.Auth)
		r.With().Get("/{id}", p.Get)
		r.With().Patch("/{id}", p.Update)
		r.With().Delete("/{id}", p.Delete)
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Print(err)
	}
}
