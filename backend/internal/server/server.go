package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/Paukku/ajanvarausjarjestelma/backend/internal/user/handler"
	"github.com/Paukku/ajanvarausjarjestelma/backend/internal/user/repository"
	"github.com/Paukku/ajanvarausjarjestelma/backend/internal/user/service"
	pbcommon "github.com/Paukku/ajanvarausjarjestelma/backend/pb/common"
	pbHTTP "github.com/Paukku/ajanvarausjarjestelma/backend/pb/http"
)

type ApiRegister struct {
	Path       string
	Handler    http.Handler
	AccessRole pbcommon.UserRole
}

var jwtSecret []byte

type RoleString string

var RoleStringToEnum = map[RoleString]pbcommon.UserRole{
	RoleString("Admin"):    pbcommon.UserRole_ADMIN,
	RoleString("Owner"):    pbcommon.UserRole_OWNER,
	RoleString("Employee"): pbcommon.UserRole_EMPLOYEE,
}

func RoleMiddleware(required pbcommon.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid claims", http.StatusUnauthorized)
				return
			}

			rawRole, ok := claims["role"].(string)
			if !ok {
				http.Error(w, "role must be a string", http.StatusForbidden)
				return
			}

			userRole, exists := RoleStringToEnum[RoleString(rawRole)]
			if !exists {
				http.Error(w, "unknown role", http.StatusForbidden)
				return
			}

			if userRole < required {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Run() {
	_ = godotenv.Load()
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewPostgresUserRepository(db)
	userService := service.NewUserServiceServer(userRepo)
	userHandler := handler.NewUserHandler(userService)

	converter := pbHTTP.NewBusinessCustomerAPIHTTPConverter(userHandler)
	mux := http.NewServeMux()

	_, path, createHandler := converter.CreateUserHTTPRule(nil)
	mux.Handle(path, createHandler)

	_, getUserPath, getUserHandler := converter.GetUserHTTPRule(nil)
	mux.Handle(getUserPath, getUserHandler)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
