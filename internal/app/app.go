package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"todo/internal/config"
	"todo/internal/handlers"
	"todo/internal/middlewares"
	"todo/internal/storage"
)

func Run() {
	cfg := config.NewConfig()

	db := storage.New(cfg.StoragePath)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middlewares.Auth)
	router.Handle("/", handlers.Index(db))
	router.Post("/login", handlers.Login(db))
	router.Post("/register", handlers.Register(db))
	router.Delete("/deleteTask/{task_id}", handlers.DeleteTask(db))

	server := http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	log.Println("Запуск сервера")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
