package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Paukku/ajanvarausjarjestelma/backend/internal/user/repository"
	"github.com/Paukku/ajanvarausjarjestelma/backend/internal/user/service"
	businessServices "github.com/Paukku/ajanvarausjarjestelma/backend/pb/http"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var jwtKey []byte

type contextKey string

const userIDKey contextKey = "user_id"

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	jwtKey = []byte(secret)
}

func AuthInterceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid user_id", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, int32(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}
	fmt.Println("Connected to PostgreSQL!")

	userService := &service.UserServiceServer{Repo: repository.NewPostgresUserRepository(db)}

	converter := businessServices.NewBusinessCustomerAPIHTTPConverter(userService)

	mux := http.NewServeMux()
	registerRoutes(mux, converter)

	handlerWithAuth := AuthInterceptor(mux)

	fmt.Println("🚀 HTTP server running on :8080")
	if err := http.ListenAndServe(":8080", handlerWithAuth); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func registerRoutes(mux *http.ServeMux, converter *businessServices.BusinessCustomerAPIHTTPConverter) {
	_, path, handler := converter.CreateUserHTTPRule(nil)
	mux.Handle(path, handler)
	_, path, handler = converter.GetUserHTTPRule(nil)
	mux.Handle(path, handler)
	_, path, handler = converter.GetUserByIdHTTPRule(nil)
	mux.Handle(path, handler)
}
