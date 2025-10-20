package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	common "github.com/Paukku/ajanvarausjarjestelma/backend/pb/common"
	businessServices "github.com/Paukku/ajanvarausjarjestelma/backend/pb/http"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

// Määritellään oma tyyppi context-avainelle
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

// AuthInterceptor tarkistaa JWT-tunnuksen ja lisää user_id:n kontekstiin
func AuthInterceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := int32(claims["user_id"].(float64))

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserServiceServer toteuttaa BusinessCustomerAPIHTTPService -interfacen
type UserServiceServer struct {
	DB *sql.DB
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *common.CreateUserRequest) (*common.GeneralResponse, error) {
	return &common.GeneralResponse{Success: true, Message: "User created!"}, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *common.EmptyRequest) (*common.UserList, error) {
	return &common.UserList{Users: []*common.User{}}, nil
}

func (s *UserServiceServer) GetUserById(ctx context.Context, req *common.GetUserRequest) (*common.User, error) {
	return &common.User{Uuid: "1", Name: "Test User"}, nil
}

func main() {
	// Yhdistetään PostgreSQL:ään
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}
	fmt.Println("Connected to PostgreSQL!")

	// Luo serveri ja converter
	userService := &UserServiceServer{DB: db}
	converter := businessServices.NewBusinessCustomerAPIHTTPConverter(userService)

	// Luo HTTP-mux ja rekisteröi handlerit
	mux := http.NewServeMux()

	_, path, handler := converter.CreateUserHTTPRule(nil)
	mux.Handle(path, handler)

	_, path, handler = converter.GetUserHTTPRule(nil)
	mux.Handle(path, handler)

	_, path, handler = converter.GetUserByIdHTTPRule(nil)
	mux.Handle(path, handler)

	// Lisää JWT auth-middleware
	handlerWithAuth := AuthInterceptor(mux)

	fmt.Println("HTTP server running on :8080")
	if err := http.ListenAndServe(":8080", handlerWithAuth); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
