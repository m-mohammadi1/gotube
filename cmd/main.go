package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"gotube/internal/config"
	"gotube/internal/handlers"
	"gotube/internal/middleware"
	"gotube/pkg/repository"
	"log"
	"net/http"
)

type Application struct {
	DSN  string
	Port string

	DB         *sql.DB
	Config     config.Data
	Repository repository.Repository
	Handler    handlers.Handler
	Middleware middleware.Middleware
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	data := config.New()
	// init app with config
	var app = Application{
		DSN:    "host=localhost port=5435 user=admin password=secret dbname=gotube sslmode=disable",
		Port:   data.Port,
		Config: data,
	}

	// connect to database
	conn, err := app.getDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	app.DB = conn

	// attach repos
	app.Repository = repository.New(app.DB)

	// add handlers and routes
	app.Middleware = middleware.New(app.Config, app.Repository.UserRepository)
	app.Handler = handlers.New(app.Repository, app.Config)
	router := app.getRouter()

	// server application
	log.Printf("starting app on port: %s", app.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", app.Port), router); err != nil {
		panic(err)
	}
}
