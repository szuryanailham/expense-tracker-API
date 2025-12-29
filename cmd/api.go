package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/jackc/pgx/v5"
	repo "github.com/szuryanailham/expense-tracker/internal/adapters/sqlc"
	"github.com/szuryanailham/expense-tracker/internal/categories"
	authMiddleware "github.com/szuryanailham/expense-tracker/internal/middleware"
	"github.com/szuryanailham/expense-tracker/internal/users"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// ===== CHI GLOBAL MIDDLEWARE =====
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(60 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {

		// public
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello world"))
		})

		queries := repo.New(app.db)
		userService := users.NewService(queries)
		userHandler := users.NewHandler(userService)
		categoryService := categories.NewService(queries)
		categoryHandler := categories.NewHandler(categoryService)

		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)

		// ===== PROTECTED =====
		r.Route("/", func(r chi.Router) {
			r.Use(authMiddleware.JWTAuth)
			r.Get("/categories", categoryHandler.GetCategoriesById)
			r.Post("/categories", categoryHandler.CreateCategory)
			r.Put("/categories/{id}", categoryHandler.UpdateCategory)
		})
	})

	return r
}


func ( app *application)run(h http.Handler)error {
		srv := &http.Server{
		Addr : app.config.addr,
		Handler: h,
		WriteTimeout: time.Second*30,
		ReadTimeout: time.Second*10,
		IdleTimeout: time.Minute,
	}
		log.Printf("server has started at addr %s", app.config.addr)
		return srv.ListenAndServe()
}



type application struct {
	config config
	db *pgx.Conn
}

type config struct  {
	addr string
	db dbConfig
}

type dbConfig struct  {
	dsn string
}