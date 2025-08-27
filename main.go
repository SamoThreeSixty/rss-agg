package main

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"database/sql"
	_"github.com/lib/pq"
	"github.com/samothreesixty/rss-agg/internal/db"
	"github.com/joho/godotenv"
	"github.com/go-chi/cors"
	"github.com/go-chi/chi"
)

type apiConfig struct {
	DB *db.Queries
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {	
		log.Fatal("Error loading .env file")
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("$PORT must be set")
	}

	dbUrl := os.Getenv("DB_URL")
	if portString == "" {
		log.Fatal("$DB_URL must be set")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	apiCfg := apiConfig{
		DB: db.New(conn),
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/ready", handlerReadiness)
	v1Router.Post("/ready2", handlerReadiness)
	v1Router.Get("/error", handleError)
	v1Router.Post("/users", apiCfg.handlerCreateUser)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%s", portString),
	}

	log.Println("Starting server on port", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
