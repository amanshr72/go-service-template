package main

import (
	"context"
	"database/sql"
	"go-crud2/internal/auth"
	"go-crud2/internal/health"
	"go-crud2/internal/middleware"
	"go-crud2/internal/profiling"
	"go-crud2/internal/user"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "go-crud2/docs"
)

// @title Go CRUD2 API
// @version 1.0
// @description REST + GraphQL User Service
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	_ = godotenv.Load()
	adapter := os.Getenv("DB_ADAPTER")

	var (
		repo user.Repository
		db   *sql.DB
	)

	if os.Getenv("APP_ENV") != "prod" {
		profiling.Start()
	}

	switch adapter {
	case "inmemory":
		repo = user.NewInMemoryRepository()
		log.Println("Using in-memory adapter")

	case "mongodb":
		repo = user.NewMongoRepository(connectMongo())
		log.Println("Using MongoDB adapter")

	default:
		var err error
		db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Fatalf("db open failed: %v", err)
		}
		defer func() {
			if err := db.Close(); err != nil {
				log.Printf("close db: %v", err)
			}
		}()
		if err := db.Ping(); err != nil {
			log.Fatalf("db ping failed: %v", err)
		}
		// if err := user.Migrate(db); err != nil {log.Fatalf("migration failed: %v", err)}
		repo = user.NewRepository(db)
		log.Println("Using Postgres adapter")
	}

	svc := user.NewService(repo)
	mux := http.NewServeMux()

	auth.RegisterRoutes(mux)
	health.RegisterRoutes(mux, db)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	apiMux := http.NewServeMux()
	if err := user.RegisterRoutes(apiMux, svc); err != nil {
		log.Fatalf("routes: %v", err)
	}

	// mux.Handle("/api/", middleware.Authenticate(apiMux))
	// mux.Handle("/graphql", middleware.Authenticate(apiMux))
	mux.Handle("/api/", apiMux)
	mux.Handle("/graphql", apiMux)

	server := middleware.Chain(mux, middleware.RequestID, middleware.Logger, middleware.Recovery)

	log.Println("Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}

// connectMongo returns *mongo.Collection directly — no interface{}
func connectMongo() *mongo.Collection {
	uri := os.Getenv("MONGO_URL")
	if uri == "" {
		log.Fatal("MONGO_URL not set")
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("mongo connect failed: %v", err)
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("mongo ping failed: %v", err)
	}
	return client.Database("gocrud2").Collection("users")
}
