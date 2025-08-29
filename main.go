package main

import (
	"fmt"
	"os"
	"time"
	"log"
	"net/http"
	"database/sql"
	_"github.com/lib/pq"

	"github.com/samothreesixty/rss-agg/internal/db"
	"github.com/samothreesixty/rss-agg/internal/handlers"
	"github.com/samothreesixty/rss-agg/internal/scraper"

	"github.com/joho/godotenv"
	"github.com/go-chi/cors"
	"github.com/go-chi/chi"
)

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

	apiCfg := handlers.ApiConfig{
		DB: db.New(conn),
	}

	go scraper.StartScraping(&apiCfg, 5, 1*time.Minute)

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

	v1Router.Post("/users", apiCfg.HandlerCreateUser)
	v1Router.Get("/user", apiCfg.MiddlewareAuth(apiCfg.HandlerGetUser))
	v1Router.Post("/feeds", apiCfg.MiddlewareAuth(apiCfg.HandlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.HandlerGetFeeds)
	v1Router.Post("/feed_follows", apiCfg.MiddlewareAuth(apiCfg.HandlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.MiddlewareAuth(apiCfg.HandlerGetFeedFollowsByUser))
	v1Router.Delete("/feed_follows/{feedFollowId}", apiCfg.MiddlewareAuth(apiCfg.HandlerDeleteFeedFollow))
	v1Router.Get("/posts", apiCfg.MiddlewareAuth(apiCfg.HandlerGetUserPosts))

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
