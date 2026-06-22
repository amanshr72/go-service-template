package main

import (
	"context"
	"database/sql"
	"fmt"
	"go-crud2/internal/auth"
	"go-crud2/internal/health"
	"go-crud2/internal/metrics"
	"go-crud2/internal/middleware"
	"go-crud2/internal/notification"
	"go-crud2/internal/profiling"
	"go-crud2/internal/user"
	"log"
	"net"
	"net/http"
	"os"

	_ "go-crud2/docs"
	"go-crud2/internal/user/pb"

	"google.golang.org/grpc/reflection"

	"go-crud2/internal/product"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()

	if os.Getenv("APP_ENV") != "prod" {
		profiling.Start()
	}

	repo, db := setupUserRepository()
	if db != nil {
		defer func() {
			if err := db.Close(); err != nil {
				log.Printf("close db: %v", err)
			}
		}()
	}

	notifier := notification.NewClient(os.Getenv("NOTIFICATION_API_URL"))
	userSvc := user.NewService(repo, notifier)
	mux := http.NewServeMux()

	// Public routes — no auth
	auth.RegisterRoutes(mux)
	health.RegisterRoutes(mux, db)

	// webrpc
	user.RegisterWebRPCRoutes(mux, userSvc)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.Handle("/metrics", metrics.Handler())
	mux.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./tsp/tsp-output/schema/openapi.yaml")
	})
	mux.HandleFunc("/swagger-ts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		if _, err := fmt.Fprint(w, `<!DOCTYPE html><html><head>
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist/swagger-ui.css">
		</head><body>
		<div id="swagger-ui"></div>
		<script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist/swagger-ui-bundle.js"></script>
		<script>
		SwaggerUIBundle({ url: "/openapi.yaml", dom_id: '#swagger-ui' })
		</script></body></html>`); err != nil {
			http.Error(w, "failed to render swagger ui", http.StatusInternalServerError)
			return
		}
	})

	go startGRPCServer(userSvc)

	// User routes (REST + GraphQL) — stdlib mux
	apiMux := http.NewServeMux()
	if err := user.RegisterRoutes(apiMux, userSvc); err != nil {
		log.Fatalf("routes: %v", err)
	}
	// mux.Handle("/api/", middleware.Authenticate(apiMux))
	// mux.Handle("/graphql", middleware.Authenticate(apiMux))
	mux.Handle("/api/", apiMux)
	mux.Handle("/graphql", apiMux)

	// Product routes — chi, mounted under stdlib mux
	mountProductRoutes(mux)

	server := middleware.Chain(mux, middleware.RequestID, middleware.Logger, middleware.Metrics, middleware.Recovery)

	log.Println("Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}

// setupUserRepository picks the Repository adapter based on DB_ADAPTER env var.
// Returns db as nil for inmemory/mongodb — only postgres needs a *sql.DB handle
func setupUserRepository() (user.Repository, *sql.DB) {
	switch os.Getenv("DB_ADAPTER") {
	case "inmemory":
		log.Println("Using in-memory adapter")
		return user.NewInMemoryRepository(), nil

	case "mongodb":
		log.Println("Using MongoDB adapter")
		return user.NewMongoRepository(connectMongo()), nil

	default:
		db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Fatalf("db open failed: %v", err)
		}
		if err := db.Ping(); err != nil {
			log.Fatalf("db ping failed: %v", err)
		}
		log.Println("Using Postgres adapter")
		return user.NewRepository(db), db
	}
}

func mountProductRoutes(mux *http.ServeMux) {
	productRepo := product.NewInMemoryRepository()
	productSvc := product.NewService(productRepo)
	productRouter := product.NewRouter(productSvc)

	mux.Handle("/api/v1/products/", http.StripPrefix("/api/v1/products", productRouter))
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
func startGRPCServer(svc user.Service) {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("grpc listen failed: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, user.NewGRPCServer(svc))
	reflection.Register(grpcServer)

	log.Println("gRPC server on :9090")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc serve failed: %v", err)
	}
}
