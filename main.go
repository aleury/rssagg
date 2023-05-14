package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aleury/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type application struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("PORT is not found in the environment")
	}
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("can't connect to database")
	}

	db := database.New(conn)
	app := application{DB: db}
	go startScraping(db, 10, time.Minute)

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
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerErr)
	// users
	v1Router.Get("/users", app.middlewareAuth(app.handlerGetUser))
	v1Router.Post("/users", app.handlerCreateUser)
	// posts
	v1Router.Get("/posts", app.middlewareAuth(app.handlerGetPostsForUser))
	// feeds
	v1Router.Get("/feeds", app.handlerGetFeeds)
	v1Router.Post("/feeds", app.middlewareAuth(app.handlerCreateFeed))
	// feed follows
	v1Router.Get("/feed_follows", app.middlewareAuth(app.handlerGetFeedFollows))
	v1Router.Post("/feed_follows", app.middlewareAuth(app.handlerCreateFeedFollow))
	v1Router.Delete("/feed_follows/{feedFollowID}", app.middlewareAuth(app.handlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	log.Printf("Server starting on port %v", portStr)
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portStr,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
